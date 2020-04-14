package gz_http

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	gzbolt "github.com/florinutz/gz-boltdb"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

// GobResponse is here because I can't gob encode the response body (which is a Reader) or the TLS field
type GobResponse struct {
	Status           string
	StatusCode       int
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             []byte
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Trailer          http.Header
	Request          *http.Request
}

// FromResponse fetches data from a http response. Data is supposed to be encode-able by gob,
// so readers like response.Body won't do
func (rg *GobResponse) FromResponse(r http.Response) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	rg.Status = r.Status
	rg.StatusCode = r.StatusCode
	rg.Proto = r.Proto
	rg.ProtoMajor = r.ProtoMajor
	rg.ProtoMinor = r.ProtoMinor
	rg.Header = r.Header
	rg.Body = bodyBytes
	rg.ContentLength = r.ContentLength
	rg.TransferEncoding = r.TransferEncoding
	rg.Close = r.Close
	rg.Uncompressed = r.Uncompressed
	rg.Trailer = r.Trailer
	rg.Request = r.Request

	return nil
}

// ToResponse restores a response from a GobResponse
func (rg *GobResponse) ToResponse() (r *http.Response) {
	r = new(http.Response)

	r.Status = rg.Status
	r.StatusCode = rg.StatusCode
	r.Proto = rg.Proto
	r.ProtoMajor = rg.ProtoMajor
	r.ProtoMinor = rg.ProtoMinor
	r.Header = rg.Header
	r.Body = ioutil.NopCloser(bytes.NewReader(rg.Body))
	r.ContentLength = rg.ContentLength
	r.TransferEncoding = rg.TransferEncoding
	r.Close = rg.Close
	r.Uncompressed = rg.Uncompressed
	r.Trailer = rg.Trailer
	r.Request = rg.Request

	return
}

// FetchUrls performs the requests in parallel and returns the responses (and errors).
// If not nil, the gotResponse callback is called when a response was retrieved. If it returns an error,
// no more responses will be handled.
func FetchUrls(
	requests []*http.Request,
	client http.Client,
	gotResponse func(response *http.Response) error) (responses []*http.Response, errs []error) {
	type reqRespErr struct {
		req  *http.Request
		resp *http.Response
		err  error
	}

	c := make(chan reqRespErr, len(requests))

	for i, req := range requests {
		go func(req *http.Request, output chan<- reqRespErr, id int) {
			resp, err := client.Do(req)
			output <- reqRespErr{
				req:  req,
				resp: resp,
				err:  err,
			}
		}(req, c, i)
	}

	timeout := time.Duration(3*len(requests)) * time.Second

	for i := 0; i < len(requests); i++ {
		select {
		case rre := <-c:
			if rre.err != nil {
				errs = append(errs, errors.Wrapf(rre.err, "req to '%s' failed", rre.req.URL.String()))
				continue
			}
			responses = append(responses, rre.resp)
			if gotResponse != nil {
				if err := gotResponse(rre.resp); err != nil {
					break
				}
			}
		case <-time.After(timeout):
			errs = append(errs, fmt.Errorf("timeout after %s", timeout))
		}
	}

	return
}

// DumpResponses fetches responses and dumps them into a gzipped bbolt database.
// If not nil, the gotResponse callback is called when a response was retrieved. If it returns an error,
// no more responses will be handled.
func DumpResponses(
	reqs []*http.Request,
	outputPath string,
	bucketName string,
	gzHeader *gzip.Header,
	gotResponse func(response *http.Response) error,
) (errs []error) {
	// load db from compressed outputPath of create a new tmp file for it
	db, err := gzbolt.Open(outputPath, &bbolt.Options{Timeout: 1 * time.Second}, false)
	if err != nil {
		errs = append(errs, err)
		return
	}
	defer db.Close()

	responses, es := FetchUrls(reqs, *http.DefaultClient, gotResponse)
	errs = append(errs, es...)

	defer func(responses []*http.Response) {
		for _, r := range responses {
			err := r.Body.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}(responses)

	if err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		gob.Register(http.Request{})

		for _, resp := range responses {
			key, err := GetKey(*resp.Request)
			if err != nil {
				return errors.Wrap(err, "couldn't generate key for request")
			}

			var content []byte
			content, err = EncodeResponse(resp)
			if err != nil {
				return errors.Wrap(err, "couldn't encode response to bytes")
			}

			if err = bucket.Put(key, content); err != nil {
				return errors.Wrapf(err, "couldn't save response into bucket '%s'", bucketName)
			}
		}

		return nil
	}); err != nil {
		errs = append(errs, err)
		return
	}

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		errs = append(errs, err)
		return
	}
	defer f.Close()

	err = gzbolt.Write(db, f, gzHeader)
	if err != nil {
		errs = append(errs, err)
		return
	}

	return
}

// GetKey computes the binary key to be used in bolt
func GetKey(req http.Request) (key []byte, err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(req)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't encode a key")
	}
	b := sha256.Sum256(buf.Bytes())

	return b[:], nil
}

// EncodeResponse computes the bytes to be stored in bolt for a http response
func EncodeResponse(resp *http.Response) (content []byte, err error) {
	responseForBin := GobResponse{}

	gob.Register(responseForBin)

	err = responseForBin.FromResponse(*resp)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't prepare response for converting to binary")
	}

	var encodedBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encodedBuffer)

	err = encoder.Encode(responseForBin)
	if err != nil {
		return
	}

	content = encodedBuffer.Bytes()

	return
}

// DecodeResponse decodes from bolt into a http Response
func DecodeResponse(from []byte) (*http.Response, error) {
	responseForBin := GobResponse{}

	decodedBuffer := bytes.NewReader(from)
	decoder := gob.NewDecoder(decodedBuffer)

	err := decoder.Decode(&responseForBin)
	if err != nil {
		return nil, err
	}

	return responseForBin.ToResponse(), nil
}

// GetAllResponses reads a bbolt db and looks for the existing responses
func GetAllResponses(db *bbolt.DB, bucketName []byte) (responses []*http.Response, err error) {
	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("no such bucket '%s'", bucketName)
		}
		if err := bucket.ForEach(func(k, v []byte) error {
			resp, err := DecodeResponse(v)
			if err != nil {
				return err
			}
			responses = append(responses, resp)

			return nil
		}); err != nil {
			return err
		}

		return nil
	})

	return
}

// GetResponses returns all the http.Response instances it finds in the gz
func GetResponses(path string, bucketName string) ([]*http.Response, error) {
	db, err := gzbolt.Open(path, &bbolt.Options{Timeout: 1 * time.Second}, true)
	if err != nil {
		return nil, errors.Wrap(err, "error loading db")
	}
	return GetAllResponses(db, []byte(bucketName))
}

// GetResponseFor scans the bolt db for a request's response.
// When nil, the matchFunc resolves to the key comparison.
func GetResponseFor(path string, bucketName string, request *http.Request,
	matchFunc func(r1, r2 *http.Request) bool) (*http.Response, error) {
	responses, err := GetResponses(path, bucketName)
	if err != nil {
		return nil, err
	}

	if matchFunc == nil {
		matchFunc = DefaultMatchFunc
	}

	for _, r := range responses {
		if matchFunc(r.Request, request) {
			return r, nil
		}
	}

	return nil, errors.New("no response was found")
}

func DefaultMatchFunc(r1, r2 *http.Request) bool {
	k1, err := GetKey(*r1)
	if err != nil {
		panic(err)
	}
	k2, err := GetKey(*r2)
	if err != nil {
		panic(err)
	}
	return string(k1) == string(k2)
}

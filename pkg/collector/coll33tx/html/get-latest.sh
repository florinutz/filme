#!/usr/bin/env bash
# keep these base64d so they don't show up in some web search
echo $(curl -sL https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/) | base64 > detail
echo "fetched detail html"

echo $(curl -sL https://1337x.to/popular-movies) | base64 > list
echo "fetched list html"

echo $(curl -sL https://1337x.to/search/dvdrip/1/) | base64 > list_pagination
echo "fetched list with pagination html"

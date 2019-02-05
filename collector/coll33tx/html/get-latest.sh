#!/usr/bin/env bash
# keep these base64d so they don't show up in some web search
content=$(curl -sL https://1337x.to/torrent/3570061/House-Party-1990-WEBRip-1080p-YTS-YIFY/)
echo $content | base64 > detail
echo "fetched detail html"

content=$(curl -sL https://1337x.to/popular-movies)
echo $content | base64 > list
echo "fetched list html"

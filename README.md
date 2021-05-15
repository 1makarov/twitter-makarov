```go
if err := twitterClient.FilteredStream(url.Values{
    "follow": {"123"},
    }); err != nil {
    log.Fatal()
}

for {
	line, err := twitterClient.Stream.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		break
	}
	if len(line) == 2 {
		continue
	}
	log.Println(string(line))
}
```
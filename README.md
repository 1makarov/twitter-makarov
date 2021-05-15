```go
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
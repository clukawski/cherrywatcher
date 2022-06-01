# Is Don Cherry Dead Yet?

https://isdoncherrydeadyet.ca

## Summary

The question of the century, when will the spectre of saturday night Canadian Hockey commentating finally kick the bucket, so we can talk hockey, talk about some good guys, talk about the troops.

This is a small daemon designed to use wikidata sparql queries to check if the bucket has been kicked over, and then send a request for a push notification to a local [gotify](https://gotify.net/) instance.

## Usage

```
$ ./cherrywatcher -h
Usage of ./cherrywatcher:
  -p string
    	Token used for gotify server push POST request
  -t	Test service with a different, but dead Don Cherry
```

## License

**WTFPL**

I chose this license because apparently it has some problems in some jurisdictions

![](./wtfpl-strip.jpg)

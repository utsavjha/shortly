Shortly
=======
URL Shortner in GoLang 

URL Shortner in GoLang, developed as an API exposing 2 endpoints:
1. `/ping` --> HealthCheck, returns `pong` back.
2. `/shorten` --> Shorten an array of URLs.
    > curl -X 'POST' http://shortly:8080/url -H 'content-type: application/json' -d '{"urls": ["http://abcxyz.com","http://facebook.com","http://google.com","http://amazon.com","http://bol.com"]}'
3. `/retrieve` --> Retrieve a Shortened URL Input and redirect to it's parent.
    > curl -X 'POST' http://shortly:8080/retrieve -H 'content-type: application/json' -d '{"urls": ["http://shortly:8080/a9bb4496370625c19838c7f812a1c0e946938dec","http://shortly:8080/8f2453e6d59f946e89bd77bce564cc3652a58843","http://shortly:8080/7d9f902afa6273f852bfa737e1ab6205e3e06ebd","http://shortly:8080/9308539e8d36064043c4089803d9610c31ab1ee7"]}'

## Added:
- [x] Module support
- [x] Containerised
- [x] keyDB instead of Redis - To permit offloading data to S3 as ColdStorage.
- [x] Support for distinguishing different users
- [x] Concurrent GoRoutines reading and writing from Database, using Channels
- [x] Same inputType struct for posting URLs as input to Shortly 


## ToDO:
- [] Generic Type support for accepting any argument to shorten: Array or Single Input
- [] Can I do authentication? 
- [] Add S3 (`localstack` dockerimage helpful.)
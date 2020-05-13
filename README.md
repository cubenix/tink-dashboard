# tink-wizard

Clone the repository with:
```
$ git clone https://github.com/gauravgahlot/tink-wizard.git && cd tink-wizard
```

The `Makefile` is all set to get you started quickly.
In order to host the server securely, please update the CA and server configuration in respective files under `/tls`.
However, the definitions are preset to give you a quick start. 
Once you have updated the `/tls` definitions, all you need is to run:
```
$ make run
```
You can now access tink-wizard at [https://localhost:7676](https://localhost:7676).

You can also start the server in an insecure mode by executing the following:
```
$ make run-insecure
```
You can now access tink-wizard at [http://localhost:7676](http://localhost:7676).

Please refer to the `Makefile`, for all that is happening behind the scenes.

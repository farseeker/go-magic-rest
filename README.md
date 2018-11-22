<h1>Golang REST-to-Magic (eDeveloper, uniPaas, XPA) proxy.</h1>

Very much just a proof of concept prototype at the moment.

Required IIS features:

- [Application Request Routing](https://www.iis.net/downloads/microsoft/application-request-routing)
  - After installation, enable Application Request Routing at the base level of your IIS server
  - This will also install [URL Rewrite 2.1](https://www.iis.net/downloads/microsoft/url-rewrite) or later if not already installed

Define the following server variables in URL Rewrite:

 - `HTTP_X_Forwarded_Appname`
 - `HTTP_X_Forwarded_Host`
 - `HTTP_X_Forwarded_Proto`
 - `HTTP_X_Scripts_Path`

 Bind the `iis` folder as an IIS website. In the IIS rewrite rule for the test site, set `HTTP_X_Forwarded_Appname` under Server Variables to be the name of the APPNAME you want
 to forward requests to (e.g. `MGRestTest`)

When making requests, the path must be exactly `/api/PRGNAME` (e.g. `/api/RESTTest`).

This will proxy a request to a URL of `http://localhost/uniScripts/mgrqispi.dll?APPNAME=MGRestTest&PRGNAME=RESTTest`

The body of the request should be a JSON blob of key/value pairs. The key/value pairs will be converted
into form parameters on the submission.
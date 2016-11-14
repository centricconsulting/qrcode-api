# Centric QR Code Generator

## Generating a QR Code
The primary focus of this set of services is to generate a QR Code that you can use as part of a web page,
mobile app or for printing out.

* `POST https://qrcode.centri.cc:3022/encode` - put the desired URL in the request body like this:

    {"url":"http://pepsi.com","size":200}.  

The `size` can range from 25x25 to 800x800.  It's optional and will default to `200` if not supplied.  It returns a PNG graphic file of the bar code.

## Other Functions
These are convenience function that can be used as part of a status page showing this component is currently operating or whatever.

* `GET https://qrcode.centri.cc:3022/ping` - Returns a JSON object - {"payload":"PONG"}
* `GET https://qrcode.centri.cc:3022/ver` - Returns a JSON object - {"payload":"1.0.3"}

NOTE: This application is currently running a self-signed certificate so in order to access it, you will probably have to turn off the certificate validation checking on whatever client you are using.

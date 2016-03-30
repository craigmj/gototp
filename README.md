gototp
======

Go Time-Based One Time Password implemention of RFC6238 (http://tools.ietf.org/html/rfc6238)

Based on pyotp: https://github.com/nathforge/pyotp

USAGE
=====
gototp is easy to use:

1. Clone the repository and build the package
```
git clone <repo-name>
cd <repo-name>
go build
```

2. Create another package and include this package in your project

3. Generate a random secret to store for a user:
````
    // Do this somewhere early, and only once
    rnd := rand.New(rand.NewSource(time.Now().Unix())
    
    // A secret length of 10 gives a 16 character secret key
    secret := gototp.RandomSecret(10, rand.New(rand.NewSource(time.Now().Unix())))
````

4.  Create the OTP object:
````
    // Create the OTP
    otp, err := gototp.New(secret)
    if nil!=err {
      panic(err)
  	}
````

5.  Find the current TOTP code:
````
    code := otp.Now()

    // Or find the previous code and the next code
    previousCode := otp.FromNow(-1)
    nextCode := otp.FromNow(1)
````

6.  Generate a Google Charts URL for a QR Code of the Secret, with a label
````
    // Google Charts URL to display a QR Code, of width (and height) 300px
    url := otp.QRCodeGoogleChartsUrl("My Own TOTP", 300)
````

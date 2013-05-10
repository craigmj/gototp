gototp
======

Go Time-Based One Time Password implemention of RFC6238 (http://tools.ietf.org/html/rfc6238)

Based on pyotp: https://github.com/nathforge/pyotp

USAGE
=====
gototp is easy to use:

1. Generate a random secret to store for a user:

    // Do this somewhere early, and only once
    rnd := rand.New(rand.NewSource(time.Now().Unix())
    
    // A secret length of 10 gives a 16 character secret key
    secret := gototp.RandomSecret(10, rand.New(rand.NewSource(time.Now().Unix())))
    
2.  Create the OTP object:

    // Create the OTP
    otp, err := gototp.New(secret)
    if nil!=err {
      panic(err)
  	}

3.  Find the current TOTP code:

    code := otp.Now()

    // Or find the previous code and the next code
    previousCode := otp.FromNow(-1)
    nextCode := otp.FromNow(1)

4.  Generate a Google Charts URL for a QR Code of the Secret, with a label

    // Google Charts URL to display a QR Code, of width (and height) 300px
    url := otp.QRCodeGoogleChartsUrl("My Own TOTP", 300)

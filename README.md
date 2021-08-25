# About the Project
Backend server for ELECT Web App written in GO(golang) with PostgreSQL Database, Microsoft Azure for Storage, and SendPulse as SMTP Service Provider.

It is a RESTful API created using:
* Gin Web Framework: [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* GORM (v1): [https://github.com/jinzhu/gorm](https://github.com/jinzhu/gorm)
* Microsoft Azure Blob Storage SDK: [https://github.com/Azure/azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go)
* Gin-Swagger (For API documentation): [https://github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
* QOR Admin (For admin interface): [https://github.com/qor/admin](https://github.com/qor/admin)
* JWT-Go (For implementing JWTs & Refresh Token): [https://github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)
* SecureCookie (For encrypting cookies): [https://github.com/gorilla/securecookie](https://github.com/gorilla/securecookie)
* Casbin (For RBAC authorization): [https://github.com/casbin/casbin](https://github.com/casbin/casbin)
* Excelize (For parsing excel files): [https://github.com/qax-os/excelize](https://github.com/qax-os/excelize)
* Go-OTP (For generating Time-based OTP): [github.com/hgfischer/go-otp](github.com/hgfischer/go-otp)
* Gomail (For sending emails): [https://github.com/go-gomail/gomail](https://github.com/go-gomail/gomail)


# API Documentation
<a href="https://e1ect.herokuapp.com/docs" target="_blank">https://e1ect.herokuapp.com/docs</a>

# ELECT
Elect is a web application for conducting college elections(specifically designed for St. Aloysius College(Autonomous), Mangalore).

# Objective of the Project
* To provide an easy to use web application for maintaining and conducting college elections.
* Secure the voting process by providing 2-factor user authentication, anonymity of the votes.

# Types of Users
0. Student
1. Admin:
    - Class Guides
    - EC/CC Club Heads
    - Student Council Coordinators
2. Super-Admin:
    - IT Administrators

# Features
* 2-Factor Authentication is implemented. OTP is be sent to the user's email during login for verification.
* Votes are completely anonymous. There will be no connection between the voter and the candidate after voting, not even in the database.
* Users are restricted to a single concurrent session(i.e, a user cannot be logged in from 2 devices at the same time).
* It has a modular approach, elections among students can be conducted for any purpose.
* Elections can be gender-specific(results will be generated separately for Male, Female and Others).

# Fallbacks
* Testing is not done :(

# Screenshots of the Web App
<img src="https://i.ibb.co/LhY3BNc/Elect1-1.png" width="600" /><br>
<img src="https://i.ibb.co/p4DLdwz/Elect2.png" width="600" />
<img src="https://i.ibb.co/52DXN2Q/image.png" height="320" /><br>
<img src="https://i.ibb.co/sJFqw6q/Elect3.png" width="600" /><br>
<img src="https://i.ibb.co/2tvFr9W/Elect4.png" width="600" /><br>
<img src="https://i.ibb.co/Sr2T2Wn/Elect10.png" width="600" /><br>
<img src="https://i.ibb.co/cb7GrFm/Elect5.png" width="600" /><br>
<img src="https://i.ibb.co/FV0NvcG/Elect6.png" width="600" /><br>
<img src="https://i.ibb.co/NrzqT6w/Elect7-1.png" width="600" /><br>
<img src="https://i.ibb.co/Xb2qmn0/Elect8.png" width="600" /><br>
<img src="https://i.ibb.co/X2qtLXf/Elect9.png" width="600" /><br>
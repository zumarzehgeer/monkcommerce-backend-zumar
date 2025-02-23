# Monk Commerce Backend Assessment

## Description
This project is a coupon management API for an E-commerce Website. It aims to assess my skills as a backend developer. 
With this tool, you can:
- POST /coupons: Create a new coupon.
- GET /coupons: Retrieve all coupons.
- GET /coupons/{id}: Retrieve a specific coupon by its ID.
- PUT /coupons/{id}: Update a specific coupon by its ID.
- DELETE /coupons/{id}: Delete a specific coupon by its ID.
- POST /applicable-coupons: Fetch all applicable coupons for a given cart and calculate the total discount that will be applied by each coupon.
- POST /apply-coupon/{id}: Apply a specific coupon to the cart and return the updated cart with discounted prices for each item.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)

## Installation
1. Clone the repository:
   ~ git clone https://github.com/zumarzehgeer/monkcommerce-backend-zumar.git
2. Navigate to the project repository
   ~ cd repo
3. Install the project dependencies
   ~ go mod download
   ~ go mod tidy
4. Start the project. You can either build that or run it directly
   ~ go run main.go (To run the project directly)
   ~ go build -o *name_you_like* main.go

## Usage
POSTMAN collection file has been included. Please Import it into the postman collection and then use it after running the project locally.

## Contact
If you have any questions, feel free to reach out:
- Email: zumarzehgeer007@gmail.com
- LinkedIn: https://www.linkedin.com/in/zumarzehgeer/
- GitHub: https://github.com/zumarzehgeer

# Currency Web Service

Welcome to the Currency Web Service! This project provides a simple and efficient way to retrieve the current USD to UAH exchange rate, subscribe to email notifications, and receive daily updates directly to your inbox.

## Features

1. **Get Current USD to UAH Rate**  
   Retrieve the latest exchange rate between USD and UAH.
   - **Endpoint:** `GET /rate`
   - **Response:** Number representing the current exchange rate.

2. **Subscribe to USD to UAH Rate Notifications**  
   Subscribe your email to receive daily updates on the USD to UAH exchange rate.
   - **Endpoint:** `POST /subscribe`
   - **Request Body:** JSON object containing your email address.
   - **Response:** Confirmation message upon successful subscription.

3. **Daily Email Notifications**  
   Receive the current USD to UAH exchange rate in your email inbox every midnight.

## Getting Started

### Prerequisites

Ensure you have Go installed on your system. You can download and install Go from [here](https://golang.org/dl/).

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/danitze/currency_web_service.git
2. Navigate to project directory:
   cd currency-web-service
3. Install the necessary dependencies
   go mod tidy
4. Run project

### Contact
If you have any questions or need further assistance, feel free [to contact me](mailto:danandryeyev@gmail.com).




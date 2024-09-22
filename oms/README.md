# Cryptocurrency Order Management System

## 1. Overview

The Order Management System (OMS) will be responsible for managing the lifecycle of orders across multiple exchanges (Binance, OKX, Bybit), ensuring that orders are placed, tracked, and updated correctly. It will also integrate risk management and portfolio tracking systems to ensure orders are compliant with risk policies. 



## 2. Core Requirements

### 2.1 Functional Requirements
1. Order Placement
- Support for different types of orders (market, limit, stop-limit, stop-loss, take-profit)
- Handle multiple exchanges (Binance, OKX, Bybit) with exchange-specific parameters (e.g. different order types or limits)
2. Order State Management
- Manage the lifecycle of orders: from **Pending -> Submitted -> Partially Filled -> Filled / Canceled / Rejected**
- Monitor order status changes and update order records accordingly
3. Error Handling & Retries
- Manage API errors gracefully (e.g. rate limits, timeouts) and retry if necessary
- Handle exchange-specific errors, order rejections, or cancellations
4. Order Matching 
   - For internal trading, support order matching to find opposing orders (e.g. buy/sell) within the system before sending them to the exchange (useful for internal liquidity management or arbitrage strategies)
5. Order Validation & Risk Checks
   - Validate orders based on exchange rules (e.g. minimum order size, price limits)
   - Perform risk checks before submitting (position limits, stop-loss, leverage limits)
6. Order Cancellations
   - Allow order cancellation (both partial and full) by the user or automatically (e.g. when a stop-loss is triggered)
7. Order History and Logging
   - Store order logs for auditing, including timestamps, exchange-specific order IDs, status changes, and result of execution
8. Concurrency
   - Handle multiple orders concurrently, potentially for different markets on different exchanges


### 2.2 Non-Functional Requirements

1. Scalability
   - Handle a high volume of orders from multiple users and exchanges without degradation in performance
2. Reliability
   - Ensure high availability and reliability for placing and tracking orders, as missed orders can result in financial loss
3. Performance
   - Maintain low-latency responses when placing orders to exchanges, especially for time-sensitive strategies
4. Security
   - Ensure secure handling of exchange API keys, use encryption for sensitive information, and employ rate-limiting to avoid being banned by the exchanges



## 3. Architecture

### 3.1 High-Level Components
```
+------------------------------------------+
|             User Interface               |
+------------------------------------------+
             |                   |
   +-------------------+      +----------------+
   | Order Input (UI)   |      |   Order Status |
   +-------------------+      +----------------+
           |                      |
+------------------------------------------+
|        Order Management System (OMS)     |
+------------------------------------------+-------+
| Order Validator | Trade Executor | Order Tracker |
+-----------------+----------------+---------------+
           |              |               |
   +----------------------------+    +------------------+
   |     Risk Management Module  |    |   Exchange APIs   |
   +----------------------------+    +------------------+
           |                                 |
   +----------------------------+    +------------------+
   | Portfolio Manager / P&L Calc |    | Binance, OKX,    |
   +----------------------------+    |    Bybit          |
                                      +------------------+
```


### 3.2 Components Breakdown

#### 3.2.1 Order Input (UI/API Gateway)
- Description: The frontend where users interact with the system, entering orders or using an API to submit programmatic orders
- Key Functions:
  - Submit buy/sell orders (market, limit, stop-loss, etc)
  - View order status, order history, and positions

#### 3.2.2 Order Validator
- Description: Validates order parameters before submission
- Key Functions:
  - Check that order size meets exchange requirements
  - Verify price limits for limit orders (against exchange-defined limits)
  - Perform risk checks (e.g. maximum exposure, leverage limits)
  - Check available balance to ensure enough funds to execute the order

#### 3.2.3 Trade Executor
- Description: The core component responsible for sending orders to the exchanges
- Key Functions:
  - Use exchange APIs (via a library like `ccxt`) to place orders 

#### 3.2.4 Order Tracker
- Description: Tracks the status of each order by monitoring updates from exchanges
- Key Functions:
  - Poll or subscribe to WebSocket streams from exchanges to receive updates on order status
  - Update the order state in the database (e.g. from "Pending" to "Filled" or "Partially Filled")
  - Notify users of order state changes
  - Automatically trigger actions based on status changes (e.g. place a stop-loss order once a primary order is filled)

#### 3.2.5 Risk Management Module
- Description: Enforces risk management policies to ensure trading activities adhere to user-defined risk limits
- Key Functions:
    Define rules for position sizing, leverage limits, and exposure
    Automatically reject orders that violate risk rules
    Optionally perform real-time checks (e.g. ensuring portfolio exposure remains within safe bounds before placing new orders)

#### 3.2.6 Portfolio Manager / P&L Calculator
- Description: Tracks the user's portfolio and updates positions after each order is executed
- Key Functions:
  - Keep track of open positions and real-time profit and loss (P&L)
  - Reconcile holdings after every order fill or cancellation
  - Provide the total portfolio value and exposure across exchanges

## 4. Order Lifecycle Management
Each order placed in the OMS should follow a lifecycle with the following stages:
1. Order Creation:
   - User submits an order **Created**
   - System validates it (order parameters, risk limits), if it fails move it to **Invalidated**
2. Order Submission:
   - The OMS sends the order to the respective exchange via its API
   - The system moves the order status to **Submitted**
3. Order Tracking:
   - The system tracks the order status (via WebSocket or polling)
   - Order status updates from the exchange are logged in the database
4. Order Execution:
   - If the order is **Partially Filled** or **Fully Filled**, the OMS updates the portfolio
   - If the order is **Canceled**, system logs the reason and releases the reserved funds
5. Order Closure
   - Once the order is completely filled or canceled, it is marked as **Closed**
   - Final P&L calculation is triggered

## 5. Database Design 

### 5.1 Order Table
- Table: `orders`
- Fields:
  - `order_id`: Unique ID
  - `user_id`: The user who submitted the order
  - `exchange`: Exchange on which the order was placed
  - `symbol`: The trading pair (e.g. BTC/USDT)
  - `order_type`: Market, limit, stop-limit, etc
  - `status`: Pending, Submitted, Partially Filled, Filled, Canceled, Rejected
  - `price`: The order price (for limit orders)
  - `quantity`: The order quantity
  - `filled_quantity`: Quantity filled so far
  - `timestamp`: Time when the order was placed
  - `exchange_order_id`: The unique order ID assigned by the exchange

### 5.2 Order History Table
- Table: `order_history`
- Fields:
  - `history_id`: Unique ID
  - `order_id`: Reference to the `orders` table
  - `status`: Updated order status
  - `timestamp`: Time when the status was updated
  - `message`: Any error or status message received from the exchange

### 5.3 Portfolio Table
- Table: `portfolio`
- Fields:
  - `user_id`: User identifier
  - `asset`: The asset field (e.g. BTC, USDT)
  - `quantity`: Quantity of the asset
  - `average_price`: Average acquisition price of the asset


## 6. Error Handling & Retry Mechanism

### 6.1 Retry Logic
- Exponential Backoff:
  - Implement an exponential backoff strategy when a request to an exchange API fails due to rate limits or temporary network issues
  - For example, if a request fails, wait for a short period before retrying, doubling the wait time with each subsequent failure up to a maximum limit
- Handling exchange-speicifc errors:
  - Detect and handle exchange-specific error codes (e.g. "insufficient funds", "rate limit exceeded")
  - Retry requests for transient errors (like network timeouts) but avoid retrying for permanent errors (like invalid order parameters)
- Maximum retry requests:
  - Define a maximum number of retry attempts to prevent infinite retry loops
  - After reaching the maximum retries, escalate the issue for manual intervention or notify the user

### 6.2 Order Failures
- Failure Detection:
  - Monitor responses from exchange APIs to detect order failures, such as invalid parameters, insufficient funds, or rejected orders
- Logging Failures:
  - Log detailed information about failed orders, including error codes, messages, and timestamps, in the `order_history` table for auditing and troubleshooting
- User Notifications:
  - Notify users immediately when an order fails, providing clear information about the reason for the failure
  - Implement notification channels such as email alerts, SMS, or in-app notifications based on user preferences
- Automatic Corrective Actions:
  - Parameter Adjustments:
    - For certain errors (e.g., order size too small), automatically adjust the order parameters to meet exchange requirements and attempt to resubmit the order
  - Funding Issues:
    - If an order fails due to insufficient funds, notify the user and provide options to deposit additional funds or adjust existing positions
  - Invalid Order Types:
    - If an unsupported or invalid order type is detected, reject the order early in the validation phase and inform the user to select a supported order type
# Home work 1

## Giới thiệu

Project về quy trình đặt vé xem phim một cách đơn giản. Project tập trung chủ yếu vào luồng đặt vé.

## Kiến trúc DB đơn giản

users (id, username, password, status, phone_number, ...)<br/>
movies (id, title, duration, ...)<br/>
showtimes (id, movie_id, start_time, screen_id)<br/>
seats (id, screen_id, row, col)<br/>
bookings (id, user_id, showtime_id, status [holding|confirmed|canceled], created_at)<br/>
booking_seats (booking_id, screen_id, seat_id)


## Luồng đặt vé

### 1. POST /api/v1/hold-seat
- Khi user click vào 1 ghế trên UI thì call api này để lock ghé đó trong vòng 10 phút để tránh user khác chọn trùng
- API sử dụng redis để lock

### 2. POST /api/v1/release-seat
- Check và xóa seat đã giữ của user

### 3. POST /api/v1/bookings
- Đặt vé từ những ghế đã giữ
- khi call api này user cần gửi thêm danh sách id ghế đã giữ

## Một số API khác:

### 1. POST /api/v1/users/register
### 2. POST /api/v1/users/login
### 3. POST /api/v1/showtimes
### 4. POST /api/v1/seats

## Cần cải thiện
- Các API đang thiếu phần validate data
- Thêm xuất hóa đơn
- Chưa tích hợp phần tính tiền, có thể dùng Momo/Zalopay/VNPAY/...
- Thêm các API để lấy danh sách phim (movies), suất chiếu (showtime), ghế (seats), ...
- Thêm nhiều table để quản lý thêm các mục khác như rạp nào, phòng nào, ...



# Home work 2

- Websocket đã được tích hợp vào home work 1
- Khi user call API giữ ghế, một event sẽ được gửi cho toàn bộ user đang kết nối. 
Mục đích để UI xử lý để các user chọn ghế 1 cách dễ dàng hơn
- Hệ thống socket hiện tại chỉ run được trên single instance nen cần tích hợp thêm Redis Pub/Sub để có thể run server trên nhiều instance khi lượng traffic lớn

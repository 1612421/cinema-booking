# Home work 1

## Giới thiệu

Project về quy trình đặt vé xem phim một cách đơn giản. Project tập trung chủ yếu vào luồng đặt vé.

## Hướng dẫn setup local project
### 1. Cài đặt docker và docker-compose
### 2. Run file setup.sh

- BE endpoint: localhost:8080
- FE endpoint: localhost:3000

## Kiến trúc DB đơn giản

users (id, username, password, status, phone_number, ...)<br/>
movies (id, title, duration, ...)<br/>
showtimes (id, movie_id, start_time, screen_id)<br/>
seats (id, screen_id, row, col)<br/>
bookings (id, user_id, showtime_id, status [holding|confirmed|canceled], created_at)<br/>
booking_seats (booking_id, screen_id, seat_id)

## Luồng đặt vé

### 1. POST /api/v1/hold-seat
- Khi user click vào 1 ghế trên UI thì call api này để báo cho hệ thống là user đang quan tâm và muốn giữ ghế này
- API sử dụng redis để tạm lưu lựa chọn của user

### 2. POST /api/v1/release-seat
- Check và xóa seat đã giữ của user

### 3. POST /api/v1/bookings
- Đặt vé từ những ghế đã giữ
- khi call api này user cần gửi thêm danh sách id ghế đã giữ
- Sau khi đặt xong, hệ thống sẽ gửi socket event cho toàn bộ client đang kết nối để thông báo về booking này

## Document:
API Document được tạo  trong folder docs
 hoặc truy cập vào: localhost:8080/swagger/index.html

## Cần cải thiện
- Các API đang thiếu phần validate data
- Thêm xuất hóa đơn
- Chưa tích hợp phần tính tiền, có thể dùng Momo/Zalopay/VNPAY/...
- Thêm các API để lấy danh sách phim (movies), suất chiếu (showtime), ghế (seats), ...
- Thêm nhiều table để quản lý thêm các mục khác như rạp nào, phòng nào, ...


## Demo:
Video demo: https://drive.google.com/file/d/13C4YJ8b5k_R1uIkkFnbvls6DP80OdLkN/view

## Load test
video load test: https://drive.google.com/file/d/1vDWMxpUcmbF0u9uahIG6FikO3vQOYk0E/view?usp=sharing

# Home work 2

- Websocket đã được tích hợp vào home work 1
- Khi user call API giữ ghế, một event sẽ được gửi cho toàn bộ user đang kết nối. 
Mục đích để UI xử lý để các user chọn ghế 1 cách dễ dàng hơn
- Hệ thống socket hiện tại chỉ run được trên single instance nen cần tích hợp thêm Redis Pub/Sub để có thể run server trên nhiều instance khi lượng traffic lớn
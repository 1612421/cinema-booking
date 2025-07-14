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

# Trả lời thêm các câu hỏi

## 1. Bạn thích golang ở điểm gì?. So sánh với PHP.
- Hiệu năng cao hơn nhiều so với PHP
- Hỗ trợ xử lý lập trình concurrent đơn giản và hiệu quả
- Deploy dễ dàng
- Sử dụng tốt trong các hệ thống microservice

## 2. Bạn thích giải quyết những vấn đề gì nhất trong backend?
- Optimize hệ thống, làm sao để API nhanh hơn, code tối ưu hơn
- Thiết kế các kiến trúc hệ thống
- Xử lý các vấn để về scale hệ thống

## 3. Bạn thấy vấn đề gì khó nhất trong backend?
- Xử lý các cơ chế lock để không ảnh hưởng hệ thống
- Các vấn để về scale hệ thống như: 
    - Khi  các record trong DB nhiều lên, làm sao để query nhanh và hiệu quả
    - Làm sao để tạo nhiều job (hàng ngàn job) và retry an toàn khi lỗi

## 4. Bình chọn một tình huống trong dự án mà bạn đã từng giải quyết, làm bạn tâm đắc, tự hào.
- Dự án: Hệ thống cá cược danh cho thị trường Úc (UPC Showdown).
- Vấn đề: Khi crawler quét được provider của nhà cái bất ngờ thêm 1 thể thức mới mà không báo trước.
Điều này gây ra không thể nhận dạng loại betting để tính điểm và kết quả cho các user.
Dẫn đến hệ thống tính điểm sai và trao giải không đúng trong khi có rất nhiều user đang theo dõi.
- Giải pháp mình đã triển khai: 
    - Meeting khẩn cấp để đưa ra business cho loại betting này
    - Vì code dùng các design pattern hiệu quả và vận dùng tốt các tính chất
  OOP nên mình chỉ cần code thêm 1 model cho loại betting này và overirde lại logic tính kết quả sau đó load lên
    - Với hệ thống quản lý các job hiệu quả, mình đã retry lại các job fail liên quan
    để cập nhật lại kết quả.

---
title: Tổng quan API
---

# Tổng quan API FCC

Tài liệu ngắn về các giao diện API tương thích của gateway Python. FCC chuyển yêu cầu đến provider đã cấu hình; tính năng mô hình, khả dụng, chi phí và điều khoản phụ thuộc provider. Bản Go và Oh My Molik vẫn đang thử nghiệm.

<a href="../../vi/#quickstart">← Cài đặt và hướng dẫn vận hành</a> · <a href="../../vi/#security">Bảo mật và truy cập từ xa</a>

## Trước khi gọi API

1. Khởi động FCC, cấu hình provider và mô hình trong Admin. Địa chỉ cục bộ mặc định là `http://127.0.0.1:8082`; nếu log hiển thị địa chỉ khác, dùng địa chỉ đó.
2. Liệt kê mô hình bằng `GET /v1/models?view=messages`, sau đó dùng ID có sẵn trong yêu cầu.
3. Nếu bật Proxy Authentication, gửi `Authorization: Bearer <token>` hoặc cấu hình token làm API key của client. Không đưa thông tin xác thực vào log công khai.

::: warning Truy cập từ xa
Proxy Authentication bảo vệ proxy API, không phải lớp xác thực riêng cho Admin hoặc Code. Hạn chế truy cập bằng biện pháp độc lập ở tầng mạng hoặc tunnel. Không công khai Admin hoặc Code cho mạng không tin cậy.
:::

## Các endpoint API

| Phương thức và đường dẫn | Mục đích | Streaming |
| --- | --- | --- |
| `GET /health` | Kiểm tra sức khỏe server cơ bản | Không |
| `GET /v1/models` | Liệt kê các mô hình | Không |
| `POST /v1/messages` | Messages tương thích Anthropic | Tùy chọn `stream: true` |
| `POST /v1/messages/count_tokens` | Đếm token cho Messages | Không |
| `POST /v1/chat/completions` | Chat tương thích OpenAI | Tùy chọn `stream: true` |
| `POST /v1/responses` | Responses tương thích OpenAI | Chỉ streaming |
| `POST /v1/models/{model}:generateContent` | Sinh nội dung tương thích Gemini | Không |
| `POST /v1/models/{model}:streamGenerateContent?alt=sse` | Streaming tương thích Gemini | SSE |

FCC cũng nhận các đường dẫn Gemini dưới `/v1beta/models/`. Tương thích không có nghĩa mọi trường, công cụ, provider hoặc hành vi streaming của upstream đều được hỗ trợ. Xem [ví dụ API trong repository](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-the-api-directly) và mã nguồn để biết chi tiết hiện tại.

## Yêu cầu đầu tiên

Thay mô hình mẫu bằng một ID sẵn có từ provider đã cấu hình:

```bash
curl http://127.0.0.1:8082/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

Ví dụ giả định Proxy Authentication chưa bật. Nếu đã bật, thêm bearer header như trên. Bạn cần provider hoạt động, thông tin xác thực hợp lệ, mô hình khả dụng và quyền tài khoản cần thiết.

## Lỗi và vận hành

Kiểm tra log server, ID mô hình, thông tin xác thực và giới hạn tài khoản trước; chỉ thử lại khi phù hợp. Phản hồi health cục bộ không xác nhận kết nối provider. Xem <a href="../../vi/#troubleshooting">khắc phục sự cố</a> và <a href="../../vi/#maintenance">bảo trì</a>. Trang tham khảo này là tĩnh và không gửi yêu cầu đến server của bạn.

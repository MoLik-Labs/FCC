---
title: Hướng dẫn FCC
description: Cài đặt, cấu hình và vận hành FCC an toàn.
---

[English guide](../) · [Markdown cho agent](guide.md) · [Bản lưu trước đây](../vi.html)

<a id="overview"></a>

FCC kết nối công cụ lập trình và ứng dụng API với nhà cung cấp mô hình đám mây hoặc cục bộ. Python là runtime được hỗ trợ; bản Go và Oh My Molik còn thử nghiệm.

<nav class="start-grid" aria-label="Bắt đầu tại đây">
  <a href="#quickstart"><strong>Cài đặt FCC</strong><span>Khởi chạy bản fork.</span></a>
  <a href="#clients"><strong>Kết nối công cụ</strong><span>Sử dụng agent và ứng dụng.</span></a>
  <a href="#security"><strong>Vận hành an toàn</strong><span>Hiểu cách truy cập từ xa.</span></a>
  <a href="#troubleshooting"><strong>Khắc phục sự cố</strong><span>Giải quyết lỗi thường gặp.</span></a>
</nav>

<a id="quickstart"></a>

## Cài đặt và chạy.

Khởi động cục bộ, cấu hình một nhà cung cấp rồi kiểm tra kết
nối.

**Cài fork này từ mã nguồn.** Lệnh dùng nhánh phát
triển hiện tại của MoLik-Labs/FCC, không phải bản phát hành
upstream. Python là runtime; bản Go vẫn đang thử nghiệm.

### Lấy mã của fork

Cài [Git](https://git-scm.com/downloads) và
[uv](https://docs.astral.sh/uv/getting-started/installation/). Mở terminal hoặc Windows PowerShell, chạy:

```text
git clone --branch arena/5fe6a43f-fcc https://github.com/MoLik-Labs/FCC.git
cd FCC
uv sync --locked
```

Đây là nhánh đang phát triển; sau khi tích hợp vào main,
chọn nhánh phù hợp. uv có thể cài Python 3.14.10. Bạn cần tự
cài coding agent; cách này không tạo shortcut desktop hay tự
cài công cụ tùy chọn.

### Khởi động server

Trong thư mục FCC trên Windows, macOS hoặc Linux:

```text
uv run fcc-server
```

Giữ terminal mở và dùng địa chỉ Admin trong log. Trong
terminal thứ hai, vào lại thư mục này để chạy client.

### Kết nối nhà cung cấp

Trong **Providers**, chọn nhà cung cấp, thêm
key và lưu. Trong **Model Config**, chọn mô
hình rồi bấm **Apply**. Ví dụ NVIDIA NIM dùng
`NVIDIA_NIM_API_KEY`; lấy key tại
[NVIDIA](https://build.nvidia.com/settings/api-keys). Giá, quyền truy cập và giới hạn do nhà cung cấp quyết
định.

Nếu không liệt kê được mô hình, nhập chính xác
`provider-id/model-id`. Với tài khoản
subscription, dùng
**Providers → OAuth providers**.

### Chạy agent

Cài agent riêng, rồi chạy Claude Code hoặc Codex từ thư mục
checkout; launcher có thể tự khởi động FCC:

```text
uv run fcc-claude
# hoặc
uv run fcc-codex
```

Chọn mô hình FCC trong client nếu có bộ chọn mô hình. Xem
[các launcher và tích hợp](#clients).

Chọn **Appearance** ở thanh bên Admin để dùng Forest,
Catppuccin Mocha, Nord, Dracula, Ocean, Dusk hoặc Dawn. Lựa chọn
được lưu trong trình duyệt khi có thể, không thay đổi cấu hình
server hay CLI Molik riêng.

<a id="first-request"></a>

## Gửi yêu cầu đầu tiên.

Sau khi server và nhà cung cấp sẵn sàng, liệt kê mô hình và gửi
yêu cầu cục bộ.

Ví dụ dùng cổng mặc định; nếu log hiển thị địa chỉ khác, hãy dùng
địa chỉ đó. Chọn ID mô hình có sẵn. Ví dụ NVIDIA yêu cầu nhà cung
cấp tương ứng đã được cấu hình.

```text
curl "http://127.0.0.1:8082/v1/models?view=messages"
```

```text
curl http://127.0.0.1:8082/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "nvidia_nim/nvidia/nemotron-3-ultra-550b-a55b",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

Nếu bật **Proxy Authentication**, thêm header
`Authorization: Bearer <token>`. Không để token
trong log chia sẻ hoặc lịch sử shell. Thành công phụ thuộc nhà
cung cấp, key và mô hình.

<a id="providers"></a>

## Nhà cung cấp & mô hình.

Dùng dịch vụ cloud, endpoint tương thích riêng hoặc mô hình cục
bộ.

### Cloud

Thêm key hoặc kết nối OAuth trong Providers, rồi chọn mô hình.
Có thể thêm và sắp xếp key theo số; key 1 là mặc định. Một số
nhà cung cấp cần vùng, endpoint hoặc tài khoản bổ sung.

### Custom provider

Chọn
**Providers → Custom providers → Add provider**.
Nhập base URL, định dạng API và key tùy chọn. Để trống Model
IDs để tự khám phá, hoặc nhập ID chính xác rồi lưu.

### Cục bộ

Khởi động model server trước. Dùng tiền tố
`lmstudio/`, `llamacpp/` hoặc
`ollama/` với ID chính xác. Refresh models hoặc bấm
Test sau khi khởi động.

### Định tuyến & fallback

`MODEL` là mặc định. Cấu hình tier, reasoning và
thứ tự fallback trong Model Config. Một yêu cầu thất bại có
thể tiêu thụ quota của nhiều nhà cung cấp.

Địa chỉ mặc định: LM Studio `http://localhost:1234/v1`,
llama.cpp `http://localhost:8080/v1`, Ollama
`http://localhost:11434`. Ưu tiên mô hình hỗ trợ tool
và đủ context cho agent.

[Danh sách và cấu hình chi tiết
↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#choose-a-provider)

<a id="clients"></a>

## Dùng công cụ của bạn.

Launcher terminal, editor, ứng dụng desktop, browser Code và API
client.

Với bản mã nguồn, cài agent riêng và thêm
`uv run` trước các lệnh launcher dưới đây khi đang ở
thư mục FCC.

| CLIENT | LỆNH | CLIENT | LỆNH |
| --- | --- | --- | --- |
| Claude Code | `fcc-claude` | Codex | `fcc-codex` |
| Pi | `fcc-pi` | OpenCode | `fcc-opencode` |
| Cline | `fcc-cline` | Hermes | `fcc-hermes` |
| DeepSeek Harness | `fcc-dsh` | Grok Build | `fcc-grok` |
| Muse Code | `fcc-muse` | Aider | `fcc-aider` |
| Antigravity CLI | `fcc-agy` |  |  |

### Editor và desktop

1. Cài và chạy client trên cùng máy, cùng môi trường người dùng với
   FCC.
2. Trong **Admin → Integrations**, bấm
   **Connect** ở thẻ tương ứng.
3. Reload VS Code hoặc khởi động lại ứng dụng/IDE; chọn mô hình FCC
   nếu cần. Dùng **Disconnect** để ngắt kết nối.

Phải thoát hoàn toàn Claude Desktop trước khi kết nối hoặc ngắt
kết nối. JetBrains ACP có bước thiết lập riêng; xem
[hướng dẫn từng tích hợp](https://github.com/MoLik-Labs/FCC/blob/main/README.md#integrations).

<a id="api"></a>

## Gọi API trực tiếp.

Trỏ SDK hoặc HTTP client tương thích tới FCC.

| GIAO THỨC | ENDPOINT | STREAMING |
| --- | --- | --- |
| Anthropic Messages | `POST /v1/messages` | Thêm `"stream": true` |
| OpenAI Chat Completions | `POST /v1/chat/completions` | Thêm `"stream": true` |
| OpenAI Responses | `POST /v1/responses` | Luôn stream |
| Gemini | `POST /v1/models/{model}:generateContent` | Dùng `:streamGenerateContent?alt=sse` |

Base URL mặc định của OpenAI SDK:
`http://127.0.0.1:8082/v1`; Anthropic và Gemini:
`http://127.0.0.1:8082`. Khi bật proxy authentication,
dùng token làm API key hoặc bearer token. Khả năng của từng
provider và mô hình vẫn được áp dụng.

[Ví dụ cho cả bốn API ↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-the-api-directly)

<a id="messaging"></a>

## Điều khiển qua điện thoại.

Cấu hình Discord hoặc Telegram trong
**Admin → Messaging**, rồi bấm Apply.

DISCORD

### Kết nối bot

Tạo bot trong Discord Developer Portal, bật Message Content
Intent; cấp quyền đọc, gửi, lịch sử và Manage Messages. Cấu
hình token, channel ID được phép và thư mục tuyệt đối được
phép.

TELEGRAM

### Kết nối bot

Tạo bot với BotFather, lấy user ID dạng số; cấu hình token,
user ID và thư mục tuyệt đối được phép. Với nhóm, cấp cho bot
quyền xóa tin nhắn.

Khởi động lại nếu được yêu cầu. `/stats` xem trạng
thái; `/stop` độc lập hủy công việc;
`/clear` độc lập xóa hội thoại FCC trong chat. Trả lời
một tin nhắn bằng lệnh sẽ nhắm tới yêu cầu hoặc chuỗi hội thoại
đó. Voice notes hỗ trợ NVIDIA NIM hoặc Whisper cục bộ tùy chọn;
xem
[thiết lập đầy đủ](https://github.com/MoLik-Labs/FCC/blob/main/README.md#use-from-your-phone).

<a id="workloads"></a>

VẬN HÀNH / TẢI

## Chọn mức tải phù hợp.

Bắt đầu thận trọng, theo dõi rồi điều chỉnh.

Trong **Admin → Providers → Runtime configuration**,
các nút gợi ý đặt tạm `PROVIDER_MAX_CONCURRENCY`: Light
1, Balanced 2 (mặc định), High 4. Kiểm tra giá trị rồi bấm
**Apply** để lưu. Nếu biến môi trường tiến trình khóa
trường này, nút gợi ý sẽ bị vô hiệu hóa.

### Light · 1

Điểm bắt đầu cho máy yếu, nhất là khi chạy mô hình cục bộ.
Giảm kích thước mô hình hoặc context; tránh cài voice cục bộ
nếu thiếu tài nguyên.

### Balanced · 2

Mức mặc định cho một người dùng. Theo dõi bộ nhớ, độ trễ và
lỗi.

### High · 4

Chỉ thử khi hạn mức nhà cung cấp và máy chạy mô hình còn dư
tải. Giảm xuống nếu bộ nhớ hoặc độ trễ tăng.

### Đo và điều chỉnh

Giới hạn áp dụng cho từng nhà cung cấp, không giới hạn toàn bộ
tài nguyên máy. Rate Limit và Rate Window vẫn độc lập. Mô hình
cloud chạy trên máy nhà cung cấp.

Đây là gợi ý, không tự nhận diện phần cứng hay đảm bảo tốc độ.
Chạy cục bộ phụ thuộc mô hình, context, quantization và server
inference; mức đồng thời không đổi các cấu hình đó.

<a id="security"></a>

## Vận hành an toàn.

Hiểu ranh giới tin cậy trước khi cho phép truy cập từ bên ngoài.

**Admin và Code không có lớp xác thực Admin riêng.**
Client từ xa truy cập được các endpoint này có thể thay đổi cấu
hình hoặc vận hành code session. Proxy Authentication bảo vệ proxy
API; không mặc định coi nó bảo vệ Admin hay Code. Không mở cho
mạng không tin cậy hoặc internet công cộng nếu thiếu kiểm soát
truy cập độc lập.

### Truy cập từ xa và tunnel

Cloudflare Tunnel, ngrok và Tailscale có thể cung cấp kết nối từ
xa khi được cấu hình phù hợp. Giới hạn truy cập tại lớp mạng hoặc
tunnel; chỉ có URL tunnel không đồng nghĩa với chính sách phân
quyền. Kiểm tra cấu hình của nhà cung cấp và xác minh người không
được phép không thể vào Admin hoặc Code trước khi sử dụng. HTTP từ
xa không mã hóa token và nội dung ở lớp HTTP; dùng HTTPS hoặc
tunnel riêng được bảo vệ phù hợp, kể cả với Claude Desktop.

### Bảo vệ thông tin xác thực

- Cấp quyền tối thiểu cho key, bot và tunnel; giới hạn user,
  channel và thư mục của messaging.
- Không commit token hoặc đăng công khai; xóa thông tin nhạy cảm
  khỏi log và báo cáo lỗi.
- Kiểm tra điều khoản, giá và thư mục mà agent có thể truy cập.

<a id="maintenance"></a>

## Quản lý cài đặt.

Kiểm tra phiên bản, cập nhật có chủ đích và sao lưu trước khi
gỡ.

| VIỆC CẦN LÀM | LỆNH / THAO TÁC |
| --- | --- |
| Xem phiên bản | `uv run fcc-server --version` |
| Chẩn đoán | `uv run fcc-doctor` |
| Cập nhật | Dừng FCC, chạy `git pull --ff-only` và `uv sync --locked` trong checkout. |
| Gỡ cài đặt | Dừng FCC và ứng dụng liên quan, sao lưu dữ liệu rồi xóa thư mục mã nguồn nếu không còn cần. |

Checkout chứa `.venv`; dữ liệu FCC như
`~/.fcc/` nằm riêng và không bị xóa cùng checkout.
Trong phiên bản này, `fcc-update` đã bị vô hiệu hóa và
không thay đổi dữ liệu; các bản updater cũ vẫn có thể chạy bộ cài
upstream. Script cài kế thừa cũng trỏ tới bản upstream, không phải
checkout này. Phát hành tự động của fork đang bị vô hiệu hóa cho
đến khi danh tính gói và quy trình phát hành được rà soát riêng.

[Hướng dẫn theo hệ điều hành ↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md#manage-your-installation)

<a id="troubleshooting"></a>

## Khi gặp sự cố.

Kiểm tra từng kết nối từ trong ra ngoài.

Không mở được Admin hoặc không kết nối được API

Kiểm tra server đang chạy, địa chỉ và cổng trong log, firewall
và chính sách tunnel. Không mở Admin công khai chỉ để khắc
phục lỗi.


Không thấy mô hình hoặc yêu cầu thất bại

Kiểm tra key, ID, quyền tài khoản và trạng thái provider. Với
mô hình cục bộ, chạy model server rồi refresh hoặc bấm Test.
Kiểm tra hỗ trợ tool và context.


Editor hoặc desktop client không thấy FCC

Chạy client và FCC trong cùng môi trường người dùng, kết nối
lại trong Admin và khởi động lại client. Một số tích hợp cần
thoát hoàn toàn hoặc thiết lập thêm.


Báo cáo lỗi như thế nào?

Chạy `uv run fcc-doctor` trong checkout; ghi lại
phiên bản, mô hình đang dùng, mappings, các bước tái hiện và
lỗi đầy đủ. Xóa key và dữ liệu riêng tư trước khi chia sẻ.
Issues đang tắt trên fork này; chỉ dùng tracker upstream cho
lỗi của upstream và bản phát hành của họ.

<a id="development"></a>

## Build và kiểm thử.

Python là runtime được hỗ trợ. Go vẫn là bản thử nghiệm chưa
hoàn chỉnh, không phải server thay thế.

Cài [uv](https://docs.astral.sh/uv/) rồi chạy checkout
này với dependency đã khóa:

```text
git clone --branch arena/5fe6a43f-fcc https://github.com/MoLik-Labs/FCC.git
cd FCC
uv sync --locked
uv run fcc-server
```

Chạy `./scripts/ci.sh` trên macOS/Linux hoặc
`.\scripts\ci.ps1` trên Windows. Đọc
[CONTRIBUTING.md](https://github.com/MoLik-Labs/FCC/blob/main/CONTRIBUTING.md)
trước khi tạo pull request. Nhánh phát triển vẫn còn lỗi type
check toàn repository và một số generator benchmark cũ; xem
checkpoint để biết trạng thái thực tế.

Với bản Go, đọc
[gocore/README.md](https://github.com/MoLik-Labs/FCC/blob/main/gocore/README.md)
và
[kế hoạch rewrite](https://github.com/MoLik-Labs/FCC/blob/main/docs/rewrite/00-master-plan.md). Test fixture thành công chưa chứng minh parity toàn bộ service
hoặc provider hiện tại.

ĐANG PHÁT TRIỂN / DỰ ÁN RIÊNG

### Oh My Molik

Prototype Rust và Go kết nối Chat endpoint và chỉ đọc file.
Rust có thêm MCP server stdio chỉ đọc, phạm vi giới hạn. Cả
hai chưa phải coding agent hoàn chỉnh.

[Tìm hiểu dự án →](https://github.com/MoLik-Labs/FCC/blob/main/ohmymolik/README.md)
  
[Xem bộ giao diện màu →](../molik.html)   
[Cấu hình và bảo mật MCP →](../molik.html#mcp-stdio)

<a id="pages"></a>

## Xuất bản trang này

Trang này dùng Astro Starlight. Chạy `npm ci` và `npm run build` trong `docs/site`. GitHub Actions kiểm tra pull request và xuất bản từ `main` khi GitHub Pages được cấu hình dùng Actions. Việc triển khai cần quyền repository và workflow thành công.

<a id="resources"></a>

## Khám phá thêm.

Tài liệu chi tiết và thông tin dự án.

[**English documentation** — Hướng dẫn tiếng Anh ↗](../)

[**README đầy đủ**Danh sách provider, API và tích hợp↗](https://github.com/MoLik-Labs/FCC/blob/main/README.md)

[**Đóng góp**Tiêu chuẩn và kiểm thử↗](https://github.com/MoLik-Labs/FCC/blob/main/CONTRIBUTING.md)

[**Issues upstream**Dành cho dự án và bản phát hành upstream↗](https://github.com/Alishahryar1/free-claude-code/issues)

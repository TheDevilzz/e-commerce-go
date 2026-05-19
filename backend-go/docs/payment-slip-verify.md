# เอกสารระบบ Payment และการตรวจสลิป

เอกสารนี้อธิบายระบบจ่ายเงินของโปรเจกต์นี้แบบภาษาคนทั่วไป โดยอ้างอิง API ของ NearbyShop หน้า `verify-noslip`

แหล่งอ้างอิง: https://docs.nearbyshop.xyz/slip-verify#verify-noslip

## ระบบนี้ทำอะไร

เมื่อลูกค้าอัปโหลดรูปสลิป ระบบจะทำงานตามลำดับนี้

1. รับ `order_id`, `amount`, และไฟล์รูปสลิปจากลูกค้า
2. ตรวจว่า order นี้เป็นของ user ที่ login อยู่จริง
3. ตรวจว่ายอดเงินที่ส่งมา ตรงกับยอดรวมของ order
4. เซฟรูปสลิปไว้ในโฟลเดอร์ `slips`
5. อ่าน QR Code จากรูปสลิป
6. ส่งข้อมูล QR ที่อ่านได้ไปตรวจที่ NearbyShop API
7. ถ้า NearbyShop ตอบว่า `success` จึงบันทึก payment
8. เปลี่ยนสถานะ order เป็น `paid`

แนวคิดสำคัญคือ ระบบจะไม่เชื่อแค่รูปสลิปอย่างเดียว แต่จะอ่าน QR ในสลิปแล้วให้ NearbyShop ตรวจว่าเป็นรายการโอนเงินจริง

## API ที่ฝั่ง Backend เปิดให้ใช้

Endpoint ของโปรเจกต์เรา:

```text
POST /api/payment
```

ต้อง login ก่อน เพราะ route นี้อยู่หลัง middleware ตรวจ token

ส่งข้อมูลแบบ `multipart/form-data`

| ชื่อ field | ตัวอย่าง | ความหมาย |
| --- | --- | --- |
| `order_id` | `12` | เลข order ที่ต้องการจ่ายเงิน |
| `amount` | `500.00` | ยอดเงินที่ต้องตรงกับ order |
| `image` | `slip.jpg` | รูปสลิปที่มี QR Code |

ตัวอย่าง cURL:

```bash
curl -X POST "http://localhost:3001/api/payment" \
  -H "Authorization: Bearer USER_LOGIN_TOKEN" \
  -F "order_id=12" \
  -F "amount=500.00" \
  -F "image=@slip.jpg"
```

ถ้าสำเร็จจะได้ประมาณนี้:

```json
{
  "message": "payment verified",
  "url": "/slips/1770000000000000000.jpg",
  "data": {
    "status": "success",
    "data": {
      "trans_id": "016064182241APM04659",
      "type": "BANK_TRANSFER",
      "amount": 500,
      "date_time": "2026-03-05T18:22:00.000Z",
      "sender": {
        "name": "นาย สมชาย ใจดี",
        "bank": "SCB",
        "account_no": "xxx-x-xx123-4"
      }
    }
  },
  "payment": {
    "id": 1,
    "user_id": 3,
    "order_id": 12,
    "ref": "016064182241APM04659",
    "amount": 500,
    "slip_url": "/slips/1770000000000000000.jpg",
    "status": "verified"
  }
}
```

## API ของ NearbyShop ที่ใช้

ระบบเรียก:

```text
POST https://api.nearbyshop.xyz/slipVerify/noSlip
```

ส่งข้อมูลแบบฟอร์ม:

| ชื่อ field | ความหมาย |
| --- | --- |
| `qr_payload` | ข้อความดิบที่อ่านจาก QR Code ในสลิป |
| `amount` | ยอดเงินที่ต้องการตรวจ เช่น `500.00` |

Header:

```text
Authorization: Bearer NEARBYSHOP_TOKEN
```

ค่า `NEARBYSHOP_TOKEN` อยู่ในไฟล์ `.env`

## ไฟล์ที่แก้

### `internal/handler/payment_handler.go`

ไฟล์นี้เป็นหัวใจของระบบ payment มีหน้าที่รับ request, อ่าน QR, ตรวจสลิป, และบันทึกผล

โครงสร้างหลัก:

```go
type UploadSlipRequest struct {
    OrderID uint    `form:"order_id" binding:"required"`
    Amount  float64 `form:"amount" binding:"required,gt=0"`
}
```

อธิบายทีละบรรทัด:

| บรรทัด | อธิบาย |
| --- | --- |
| `type UploadSlipRequest struct {` | สร้างรูปแบบข้อมูลที่ backend คาดว่าจะได้รับจาก frontend |
| `OrderID uint ...` | ต้องส่งเลข order มาด้วย เพื่อรู้ว่าจ่ายเงินให้ order ไหน |
| `Amount float64 ...` | ต้องส่งยอดเงิน และต้องมากกว่า 0 |
| `}` | จบการกำหนดรูปแบบ request |

```go
type VerifyResponse struct {
    Status string     `json:"status"`
    Data   VerifyData `json:"data"`
}
```

อธิบาย:

| บรรทัด | อธิบาย |
| --- | --- |
| `type VerifyResponse struct {` | สร้างรูปแบบคำตอบจาก NearbyShop |
| `Status string` | เก็บคำว่า `success` หรือสถานะอื่น |
| `Data VerifyData` | เก็บรายละเอียดการโอนเงิน |
| `}` | จบโครงสร้างคำตอบ |

```go
func (h *PaymentHandler) UploadSlip(c *gin.Context) {
```

ฟังก์ชันนี้คือจุดที่ `/api/payment` เรียกเข้ามา

ลำดับการทำงานในฟังก์ชันนี้:

| ขั้น | โค้ดทำอะไร | เหตุผล |
| --- | --- | --- |
| 1 | `c.ShouldBind(&req)` | อ่าน `order_id` และ `amount` จาก form |
| 2 | `c.GetUint("userId")` | เอา user id จาก token ที่ login |
| 3 | ค้นหา order ด้วย `id` และ `user_id` | กัน user ไปจ่าย order ของคนอื่น |
| 4 | เช็ค `order.Status == "paid"` | กันจ่ายซ้ำ |
| 5 | เทียบยอด order กับ amount | กันกรอกยอดผิด |
| 6 | `c.FormFile("image")` | รับรูปสลิป |
| 7 | `h.ReadQRCode(...)` | เซฟรูปและอ่าน QR |
| 8 | เช็ค `ref` ในตาราง payment | กันใช้สลิปเดิมซ้ำ |
| 9 | สร้าง `model.Payment` | เตรียมข้อมูลที่จะบันทึกลงฐานข้อมูล |
| 10 | `tx := h.db.Begin()` | เริ่ม transaction เพื่อให้ payment กับ order อัปเดตไปด้วยกัน |
| 11 | `tx.Create(&payment)` | บันทึก payment |
| 12 | `order.Status = "paid"` | เปลี่ยน order เป็นจ่ายแล้ว |
| 13 | `tx.Save(&order)` | บันทึกสถานะ order |
| 14 | `tx.Commit()` | ยืนยันการเปลี่ยนแปลงทั้งหมด |
| 15 | `c.JSON(...)` | ส่งผลลัพธ์กลับ frontend |

```go
func (h *PaymentHandler) ReadQRCode(file *multipart.FileHeader, amount float64) (string, *VerifyResponse, error) {
```

ฟังก์ชันนี้รับรูปสลิป แล้วคืน 3 อย่าง:

| ค่าที่คืน | ความหมาย |
| --- | --- |
| `string` | path ของรูปที่เซฟไว้ เช่น `slips/xxx.jpg` |
| `*VerifyResponse` | ผลตรวจสลิปจาก NearbyShop |
| `error` | ถ้ามีปัญหา จะส่ง error กลับ |

ลำดับข้างใน:

| โค้ด | อธิบาย |
| --- | --- |
| `filepath.Ext(file.Filename)` | เอานามสกุลไฟล์ เช่น `.jpg` |
| `time.Now().UnixNano()` | ตั้งชื่อไฟล์ด้วยเวลาแบบละเอียด เพื่อลดโอกาสชื่อซ้ำ |
| `os.MkdirAll("slips", os.ModePerm)` | สร้างโฟลเดอร์ `slips` ถ้ายังไม่มี |
| `file.Open()` | เปิดไฟล์ที่ user อัปโหลด |
| `os.Create(path)` | สร้างไฟล์ใหม่ในเครื่อง |
| `dst.ReadFrom(src)` | คัดลอกรูปจาก request ไปเก็บในเครื่อง |
| `image.Decode(QR)` | แปลงไฟล์รูปให้ Go อ่านได้ |
| `gozxing.NewBinaryBitmapFromImage(img)` | เตรียมรูปให้ library อ่าน QR |
| `qrReader.Decode(bitmap, nil)` | อ่านข้อความจาก QR Code |
| `h.CheckPayment(result.GetText(), amount)` | เอาข้อความ QR ไปตรวจที่ NearbyShop |

```go
func (h *PaymentHandler) CheckPayment(qrPayload string, amount float64) (*VerifyResponse, error) {
```

ฟังก์ชันนี้คือส่วนที่เรียก NearbyShop

ลำดับข้างใน:

| โค้ด | อธิบาย |
| --- | --- |
| `url := "https://api.nearbyshop.xyz/slipVerify/noSlip"` | กำหนด endpoint ตามเอกสาร NearbyShop |
| `os.Getenv("NEARBYSHOP_TOKEN")` | อ่าน token จาก `.env` |
| `multipart.NewWriter(body)` | เตรียม body แบบ form-data |
| `writer.WriteField("qr_payload", qrPayload)` | ใส่ QR payload ที่อ่านจากรูป |
| `writer.WriteField("amount", fmt.Sprintf("%.2f", amount))` | ใส่ยอดเงินเป็นทศนิยม 2 ตำแหน่ง |
| `http.NewRequest(http.MethodPost, url, body)` | สร้าง request แบบ POST |
| `req.Header.Set("Authorization", "Bearer "+token)` | ใส่ token เพื่อให้ NearbyShop ยอมรับ |
| `req.Header.Set("Content-Type", writer.FormDataContentType())` | บอกว่า body เป็น form-data |
| `client.Do(req)` | ส่ง request ไป NearbyShop |
| `io.ReadAll(resp.Body)` | อ่านคำตอบกลับมา |
| `json.Unmarshal(data, &verifyResp)` | แปลง JSON เป็น struct ของ Go |
| `verifyResp.Status != "success"` | ถ้าไม่ success ให้ถือว่าจ่ายไม่ผ่าน |
| `return &verifyResp, nil` | ส่งผลตรวจกลับไปใช้ต่อ |

### `internal/model/payment.go`

ไฟล์นี้คือรูปแบบตาราง payment ใน database

| Field | ความหมาย |
| --- | --- |
| `ID` | เลข payment |
| `UserID` | user ที่จ่ายเงิน |
| `OrderID` | order ที่ถูกจ่าย |
| `Ref` | เลขอ้างอิงธุรกรรมจากสลิป ใช้กันสลิปซ้ำ |
| `Amount` | จำนวนเงิน |
| `SlipURL` | path รูปสลิป |
| `Status` | สถานะ payment เช่น `verified` |
| `CreatedAt` | เวลาที่สร้าง |
| `UpdatedAt` | เวลาที่แก้ล่าสุด |

### `internal/router/router.go`

เพิ่มบรรทัดนี้:

```go
r.Static("/slips", "./slips")
```

ความหมายคือ ถ้ามีไฟล์อยู่ที่ `slips/a.jpg` จะเปิดผ่าน URL ได้เป็น:

```text
/slips/a.jpg
```

## ทำไมต้องใช้ transaction

ตอนบันทึก payment ระบบต้องทำ 2 อย่างพร้อมกัน

1. เพิ่ม record ในตาราง payment
2. เปลี่ยน order เป็น `paid`

ถ้าทำอย่างแรกสำเร็จ แต่อย่างที่สองล้มเหลว ข้อมูลจะไม่ตรงกัน ดังนั้นจึงใช้ transaction เพื่อให้สำเร็จพร้อมกัน หรือยกเลิกพร้อมกัน

## Error ที่อาจเจอ

| ข้อความ | สาเหตุ |
| --- | --- |
| `order not found` | order ไม่มีจริง หรือไม่ใช่ของ user นี้ |
| `order is already paid` | order นี้เคยจ่ายแล้ว |
| `amount does not match order total` | ยอดที่ส่งมาไม่ตรงกับยอด order |
| `image is required` | ไม่ได้แนบรูป field ชื่อ `image` |
| `failed to read QR code from image` | รูปไม่มี QR หรืออ่าน QR ไม่ได้ |
| `NEARBYSHOP_TOKEN is not configured` | ยังไม่ได้ตั้ง token ใน `.env` |
| `this slip has already been used` | เลขธุรกรรมนี้เคยถูกบันทึกแล้ว |

## สรุปแบบสั้นมาก

ระบบนี้ไม่ใช่แค่รับรูปสลิปแล้วบอกผ่าน แต่จะตรวจว่า:

1. order เป็นของ user จริง
2. ยอดเงินตรงกับ order
3. QR ในสลิปตรวจผ่าน NearbyShop
4. ref จากสลิปยังไม่เคยถูกใช้
5. ถ้าทุกอย่างผ่าน จึงบันทึก payment และเปลี่ยน order เป็น `paid`

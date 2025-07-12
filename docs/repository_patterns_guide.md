# Repository Patterns Guide: ToDomain, CreateModelFromDomain & ReadModel

## 📋 สารบัญ

1. [ภาพรวมการแยกโมเดล](#1-ภาพรวมการแยกโมเดล)
2. [Database Model vs Domain Model vs Read Model](#2-database-model-vs-domain-model-vs-read-model)
3. [การ Implement ToDomain และ CreateModelFromDomain](#3-การ-implement-todomain-และ-createmodelfromdomain)
4. [การใช้งาน Read Models](#4-การใช้งาน-read-models)
5. [ตัวอย่างการใช้งานใน Repository](#5-ตัวอย่างการใช้งานใน-repository)
6. [Best Practices](#6-best-practices)

---

## 1. ภาพรวมการแยกโมเดล

ใน Clean Architecture เราแยกโมเดลออกเป็น 3 ประเภทหลัก:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Database Model │◄──►│   Domain Model  │◄──►│   Read Model    │
│   (GORM Struct) │    │ (Business Logic)│    │ (Query Results) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        │                        │                        │
    ┌───▼────┐               ┌───▼────┐               ┌───▼────┐
    │Database│               │Business│               │   API  │
    │ Layer  │               │ Logic  │               │Response│
    └────────┘               └────────┘               └────────┘
```

**เหตุผลที่ต้องแยก:**

- **Database Model**: เพื่อ ORM mapping และ database constraints
- **Domain Model**: เพื่อ business logic และ validation
- **Read Model**: เพื่อ performance และ API response ที่เหมาะสม

---

## 2. Database Model vs Domain Model vs Read Model

### Database Model (GORM Struct)

```go
type UserModel struct {
    ID        uint64 `gorm:"primaryKey;autoIncrement"`
    Username  string `gorm:"type:varchar(50);uniqueIndex;not null"`
    Email     string `gorm:"type:varchar(255);uniqueIndex;not null"`
    FirstName string `gorm:"type:varchar(100);not null"`
    LastName  string `gorm:"type:varchar(100);not null"`
    Status    string `gorm:"type:varchar(20);default:'active'"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**จุดประสงค์:**

- ✅ ORM mapping
- ✅ Database constraints
- ✅ Table relationships
- ❌ ไม่มี business logic
- ❌ ไม่มี validation

### Domain Model (Business Entity)

```go
type User struct {
    id        ID                // Value Object with validation
    username  NonEmptyString    // Value Object with validation
    email     NonEmptyString    // Value Object with validation
    firstName NonEmptyString    // Value Object with validation
    lastName  NonEmptyString    // Value Object with validation
    status    Status            // Enum with business rules
    createdAt CreatedAt         // Value Object with validation
    updatedAt UpdatedAt         // Value Object with validation
}

// Business methods
func (u *User) IsActive() bool { ... }
func (u *User) Activate() error { ... }
func (u *User) UpdateProfile(firstName, lastName string) error { ... }
```

**จุดประสงค์:**

- ✅ Business logic
- ✅ Data validation
- ✅ Encapsulation
- ✅ Immutability
- ❌ ไม่รู้จัก database
- ❌ ไม่รู้จัก JSON serialization

### Read Model (Query Results)

```go
type UserListItem struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FullName  string    `json:"full_name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type UserProfile struct {
    ID             int64      `json:"id"`
    Username       string     `json:"username"`
    Email          string     `json:"email"`
    FirstName      string     `json:"first_name"`
    LastName       string     `json:"last_name"`
    FullName       string     `json:"full_name"`
    Status         string     `json:"status"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
    ProfilePicture *string    `json:"profile_picture,omitempty"`
    PostCount      int64      `json:"post_count"`
    FollowerCount  int64      `json:"follower_count"`
}
```

**จุดประสงค์:**

- ✅ Optimized for queries
- ✅ JSON serialization
- ✅ Specific use cases (list, profile, stats)
- ✅ Join data from multiple tables
- ❌ ไม่มี business logic
- ❌ ไม่ใช้สำหรับ write operations

---

## 3. การ Implement ToDomain และ CreateModelFromDomain

### ToDomain Method

แปลง Database Model เป็น Domain Model

```go
// ToDomain แปลง Database Model เป็น Domain Entity
func (u *UserModel) ToDomain() (*domain.User, error) {
    return domain.ReconstructUser(
        int64(u.ID),
        u.Username,
        u.Email,
        u.FirstName,
        u.LastName,
        u.Status,
        u.CreatedAt,
        u.UpdatedAt,
    )
}
```

**หลักการ:**

- ใช้ `ReconstructUser` สำหรับข้อมูลที่มีอยู่แล้ว (มี ID)
- ใช้ `NewUser` สำหรับข้อมูลใหม่ (ไม่มี ID)
- Handle validation errors จาก Value Objects

### CreateModelFromDomain Functions

แปลง Domain Model เป็น Database Model

```go
// CreateModelFromDomain สำหรับ User ที่มี ID แล้ว (Update)
func CreateModelFromDomain(user *domain.User) *UserModel {
    return &UserModel{
        ID:        uint64(user.ID().Int64()),
        Username:  user.Username().String(),
        Email:     user.Email().String(),
        FirstName: user.FirstName().String(),
        LastName:  user.LastName().String(),
        Status:    user.Status().String(),
        CreatedAt: user.CreatedAt().Time(),
        UpdatedAt: user.UpdatedAt().Time(),
    }
}

// CreateNewModelFromDomain สำหรับ User ใหม่ (Create)
func CreateNewModelFromDomain(user *domain.User) *UserModel {
    return &UserModel{
        // ไม่ใส่ ID เพราะ database จะ auto-generate
        Username:  user.Username().String(),
        Email:     user.Email().String(),
        FirstName: user.FirstName().String(),
        LastName:  user.LastName().String(),
        Status:    user.Status().String(),
        CreatedAt: user.CreatedAt().Time(),
        UpdatedAt: user.UpdatedAt().Time(),
    }
}
```

**หลักการ:**

- แยกฟังก์ชันสำหรับ Create (ไม่มี ID) และ Update (มี ID)
- เรียก `.String()`, `.Int64()`, `.Time()` เพื่อดึงค่าจาก Value Objects
- ไม่ต้อง handle errors เพราะ Domain Model ถูกต้องแล้ว

---

## 4. การใช้งาน Read Models

### เมื่อไหร่ควรใช้ Read Models

**ใช้ Read Models เมื่อ:**

- Query ที่ join หลายตาราง
- การแสดงผลที่ต้องการ format พิเศษ (เช่น FullName)
- Performance optimization (select เฉพาะ fields ที่ต้องการ)
- API response ที่มี structure ต่างจาก domain

**ตัวอย่างการใช้งาน:**

```go
// ❌ ไม่ดี: ใช้ Domain Model สำหรับ list
func (r *userRepository) GetUserList(ctx context.Context) ([]*domain.User, error) {
    // ต้อง select ทุก field และแปลงเป็น domain
    // ไม่มี FullName, ไม่มี related data
}

// ✅ ดี: ใช้ Read Model สำหรับ list
func (r *userRepository) GetUserList(ctx context.Context, query *database.Query) ([]*readmodel.UserListItem, error) {
    var items []*readmodel.UserListItem

    db := r.db.WithContext(ctx).Model(&UserModel{}).
        Select("id, username, email, CONCAT(first_name, ' ', last_name) as full_name, status, created_at")

    // Apply filters, sorting, pagination
    if query != nil {
        // ... apply query conditions
    }

    err := db.Find(&items).Error
    return items, err
}
```

### Read Model Examples

```go
// สำหรับ Dashboard/Analytics
func (r *userRepository) GetUserStats(ctx context.Context) (*readmodel.UserStats, error) {
    var stats readmodel.UserStats

    // Multiple queries for different metrics
    r.db.WithContext(ctx).Model(&UserModel{}).Count(&stats.TotalUsers)
    r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", "active").Count(&stats.ActiveUsers)

    // Time-based queries
    today := time.Now().Truncate(24 * time.Hour)
    r.db.WithContext(ctx).Model(&UserModel{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday)

    return &stats, nil
}

// สำหรับ Search
func (r *userRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*readmodel.UserSearch, error) {
    var results []*readmodel.UserSearch

    searchPattern := "%" + keyword + "%"

    err := r.db.WithContext(ctx).
        Model(&UserModel{}).
        Select("id, username, email, CONCAT(first_name, ' ', last_name) as full_name, status").
        Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
            searchPattern, searchPattern, searchPattern, searchPattern).
        Limit(limit).
        Find(&results).Error

    return results, err
}
```

---

## 5. ตัวอย่างการใช้งานใน Repository

### Repository Interface Design

```go
type UserRepository interface {
    // Domain operations (สำหรับ business logic)
    Create(ctx context.Context, user *domain.User) (*domain.User, error)
    CreateInTx(tx *gorm.DB, user *domain.User) (*domain.User, error)
    GetByID(ctx context.Context, id int64) (*domain.User, error)
    GetByIDInTx(tx *gorm.DB, id int64) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    UpdateInTx(tx *gorm.DB, user *domain.User) error
    Delete(ctx context.Context, id int64) error
    DeleteInTx(tx *gorm.DB, id int64) error

    // Domain queries (คืน Domain Models)
    FindByUsername(ctx context.Context, username string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)

    // Read model operations (สำหรับ queries และ API responses)
    GetUserList(ctx context.Context, query *database.Query) ([]*readmodel.UserListItem, error)
    GetUserProfile(ctx context.Context, id int64) (*readmodel.UserProfile, error)
    GetUserStats(ctx context.Context) (*readmodel.UserStats, error)
    SearchUsers(ctx context.Context, keyword string, limit int) ([]*readmodel.UserSearch, error)
}
```

### Implementation Pattern

```go
// Domain operations
func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    return r.CreateInTx(r.db.WithContext(ctx), user)
}

func (r *userRepository) CreateInTx(tx *gorm.DB, user *domain.User) (*domain.User, error) {
    // 1. แปลง Domain → Database Model
    model := CreateNewModelFromDomain(user)

    // 2. Save to database
    if err := tx.Create(model).Error; err != nil {
        return nil, err
    }

    // 3. แปลง Database Model → Domain (พร้อม ID ใหม่)
    return model.ToDomain()
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    return r.GetByIDInTx(r.db.WithContext(ctx), id)
}

func (r *userRepository) GetByIDInTx(tx *gorm.DB, id int64) (*domain.User, error) {
    // 1. Query Database Model
    var model UserModel
    err := tx.First(&model, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    // 2. แปลง Database Model → Domain
    return model.ToDomain()
}

// Read model operations
func (r *userRepository) GetUserList(ctx context.Context, query *database.Query) ([]*readmodel.UserListItem, error) {
    var items []*readmodel.UserListItem

    // Query โดยตรงไปยัง Read Model (ไม่ผ่าน Domain)
    db := r.db.WithContext(ctx).Model(&UserModel{}).
        Select("id, username, email, CONCAT(first_name, ' ', last_name) as full_name, status, created_at")

    // Apply query conditions
    if query != nil {
        // ... apply filters, sorting, pagination
    }

    err := db.Find(&items).Error
    return items, err
}
```

---

## 6. Best Practices

### ✅ DO's

1. **แยก Domain และ Read operations ใน Repository**

   ```go
   // Domain operations คืน domain.User
   GetByID(ctx context.Context, id int64) (*domain.User, error)

   // Read operations คืน readmodel
   GetUserProfile(ctx context.Context, id int64) (*readmodel.UserProfile, error)
   ```

2. **ใช้ Read Models สำหรับ complex queries**

   ```go
   // ✅ ดี: Query optimized สำหรับ use case
   SELECT id, username, CONCAT(first_name, ' ', last_name) as full_name,
          (SELECT COUNT(*) FROM posts WHERE user_id = users.id) as post_count
   FROM users WHERE status = 'active'
   ```

3. **Handle errors ใน ToDomain**

   ```go
   func (u *UserModel) ToDomain() (*domain.User, error) {
       user, err := domain.ReconstructUser(...)
       if err != nil {
           return nil, fmt.Errorf("invalid user data: %w", err)
       }
       return user, nil
   }
   ```

4. **แยกฟังก์ชัน Create และ Update**
   ```go
   CreateNewModelFromDomain(user)  // ไม่มี ID
   CreateModelFromDomain(user)     // มี ID
   ```

### ❌ DON'Ts

1. **อย่าใช้ Domain Models สำหรับ API responses**

   ```go
   // ❌ ไม่ดี
   type UserResponse struct {
       *domain.User  // Exposes internal structure
   }

   // ✅ ดี
   type UserResponse struct {
       ID       int64  `json:"id"`
       Username string `json:"username"`
       // ... specific fields for API
   }
   ```

2. **อย่าใส่ business logic ใน Database Models**

   ```go
   // ❌ ไม่ดี
   func (u *UserModel) Activate() {
       u.Status = "active"  // No validation!
   }

   // ✅ ดี - ใส่ใน Domain Model
   func (u *domain.User) Activate() error {
       return u.changeStatus("active")  // With validation
   }
   ```

3. **อย่าใช้ Read Models สำหรับ write operations**

   ```go
   // ❌ ไม่ดี
   func UpdateUser(profile *readmodel.UserProfile) error {
       // Read model ไม่มี business logic
   }

   // ✅ ดี
   func UpdateUser(user *domain.User) error {
       // Domain model มี validation และ business rules
   }
   ```

### 📊 Performance Tips

1. **Select เฉพาะ fields ที่ต้องการใน Read Models**

   ```go
   Select("id, username, email")  // แทนที่จะ Select("*")
   ```

2. **ใช้ Preload สำหรับ related data**

   ```go
   db.Preload("Profile").Preload("Posts", "status = ?", "published")
   ```

3. **ใช้ Raw SQL สำหรับ complex analytics**
   ```go
   db.Raw(`
       SELECT DATE(created_at) as date, COUNT(*) as count
       FROM users
       WHERE created_at >= ?
       GROUP BY DATE(created_at)
   `, time.Now().AddDate(0, 0, -30))
   ```

---

## สรุป

การแยก **Database Model**, **Domain Model**, และ **Read Model** ช่วยให้:

- 🏗️ **Clean Architecture**: แยก concerns ชัดเจน
- 🔒 **Type Safety**: Domain Models มี validation
- ⚡ **Performance**: Read Models optimized สำหรับ queries
- 🧪 **Testability**: Business logic แยกจาก database
- 🔧 **Maintainability**: เปลี่ยนแปลง database ไม่กระทบ business logic

**Pattern นี้เหมาะกับ:**

- ✅ Complex business applications
- ✅ Applications ที่ต้องการ high performance
- ✅ Teams ที่มีหลาย developers
- ✅ Long-term maintainable projects

**ไม่เหมาะกับ:**

- ❌ Simple CRUD applications
- ❌ Prototypes หรือ proof of concepts
- ❌ Projects ที่ต้องการ rapid development

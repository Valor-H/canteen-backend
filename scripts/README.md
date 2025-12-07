# 数据库初始化脚本

## 使用说明

此目录包含餐厅管理系统的数据库初始化脚本。

### 初始化数据库

1. 确保MySQL服务已启动
2. 创建数据库（如果不存在）：
   ```sql
   CREATE DATABASE IF NOT EXISTS canteen DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
   ```
3. 执行初始化脚本：
   ```bash
   mysql -u root -p canteen < scripts/init_database.sql
   ```

或者在MySQL命令行中：
```sql
USE canteen;
SOURCE scripts/init_database.sql;
```

### 脚本内容

`init_database.sql`包含以下内容：
- 创建所有必要的数据表
- 定义表结构和索引
- 示例数据（已注释，可根据需要取消注释）

### 测试数据

脚本中包含了一些测试数据的插入语句（已注释），如果需要初始化测试数据，请取消相关注释后重新执行脚本。
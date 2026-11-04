# Util工具包接口文档

本目录包含多个实用的工具包，提供各种常用功能的封装。

```md
1. 找出本项目util目录下的所有文件夹和代码文件
2. 遍历第一步找出的所有文件，整理出每个文件对外提供的接口，以及具体含义
3. 在util目录下新增一个util.md 文件
4. 在util.md 里面进行填充， 一个文件对应一个表格，表格里面的内容：里面第2步的{模块名称，接口名称，含义，请求参数，返回参数，文件地址}整 (这个表格后续会让你进行索引可以更改为更适合你 查找方法的格式)
```
## arrayutil 模块

**文件地址**: `util/arrayutil/arrayutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| arrayutil | ReverseSlice | 反转数组元素顺序 | s []T - 需要反转的切片 | 无(原地修改) | arrayutil/arrayutil.go |
| arrayutil | FilterEmptyStrings | 过滤数组中的空字符串 | arr []string - 字符串数组 | []string - 过滤后的数组 | arrayutil/arrayutil.go |
| arrayutil | DeDuplication | 数组去重 | arr []string - 需要去重的字符串数组 | []string - 去重后的数组 | arrayutil/arrayutil.go |
| arrayutil | Contains | 判断数组是否包含指定元素 | slice []T - 切片, value T - 要查找的值 | bool - 是否包含 | arrayutil/arrayutil.go |
| arrayutil | GetStructAttrArr | 将结构体数组转换为指定字段的字符串数组 | items interface{} - 结构体数组, fieldName string - 字段名(支持嵌套如"BasicInfo.OrgId") | ([]string, error) - 字段值数组和错误 | arrayutil/arrayutil.go |
| arrayutil | ArrToMapSlice | 将数组转换为map，key为指定字段，值为结构体切片 | arr []T - 结构体数组, keyName string - 作为key的字段名 | (map[string][]T, error) - 转换后的map和错误 | arrayutil/arrayutil.go |
| arrayutil | ArrToMapMapSlice | 将数组转换为双层map结构，值为结构体切片 | arr []T - 结构体数组, keyName1 string - 第一层key字段, keyName2 string - 第二层key字段 | (map[string]map[string][]T, error) - 双层map和错误 | arrayutil/arrayutil.go |
| arrayutil | ArrToMapMapStruct | 将数组转换为双层map结构，值为单个结构体 | arr []T - 结构体数组, keyName1 string - 第一层key字段, keyName2 string - 第二层key字段 | (map[string]map[string]T, error) - 双层map和错误 | arrayutil/arrayutil.go |
| arrayutil | ArrToMapStruct | 将数组转换为map，key为指定字段，值为单个结构体 | arr []T - 结构体数组, keyName string - 作为key的字段名 | (map[string]T, error) - map和错误 | arrayutil/arrayutil.go |
| arrayutil | GetMax | 获取整数数组中的最大值 | nums ...int - 可变参数整数 | int - 最大值 | arrayutil/arrayutil.go |
| arrayutil | Intersect | 计算两个字符串数组的交集 | a []string - 第一个数组, b []string - 第二个数组 | []string - 交集结果 | arrayutil/arrayutil.go |
| arrayutil | TrimSplitFilterEmpty | 去除空格后分割字符串并过滤空值 | str string - 要分割的字符串, sep string - 分隔符 | []string - 处理后的字符串数组 | arrayutil/arrayutil.go |
| arrayutil | SplitFilterEmpty | 分割字符串并过滤空值 | str string - 要分割的字符串, sep string - 分隔符 | []string - 处理后的字符串数组 | arrayutil/arrayutil.go |

## dateutil 模块

**文件地址**: `util/dateutil/dateutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| dateutil | GetBeijingLocation | 获取北京时区 | 无 | *time.Location - 北京时区对象 | dateutil/dateutil.go |
| dateutil | DefaultFormat | 格式化时间为默认格式 | t time.Time - 要格式化的时间 | string - 格式化后的字符串(2006-01-02 15:04:05) | dateutil/dateutil.go |
| dateutil | DefaultFormatYMD | 格式化时间为日期格式 | t time.Time - 要格式化的时间 | string - 格式化后的日期字符串(2006-01-02) | dateutil/dateutil.go |
| dateutil | BeijingFormat | 将时间转换为北京时间并格式化 | t time.Time - 要格式化的时间 | string - 北京时间格式化字符串(2006-01-02 15:04:05) | dateutil/dateutil.go |
| dateutil | BeijingFormatYMD | 将时间转换为北京时间并格式化为日期 | t time.Time - 要格式化的时间 | string - 北京时间日期字符串(2006-01-02) | dateutil/dateutil.go |
| dateutil | BeijingByFormat | 按指定格式将时间转换为北京时间 | t time.Time - 要格式化的时间, format string - 自定义格式 | string - 按指定格式格式化的北京时间字符串 | dateutil/dateutil.go |
| dateutil | BeijingParseYMD | 解析日期字符串为北京时间 | t string - 日期字符串(格式:2006-01-02) | (time.Time, error) - 解析后的时间和错误 | dateutil/dateutil.go |
| dateutil | BeijingDefaultParse | 解析完整时间字符串为北京时间 | t string - 时间字符串(格式:2006-01-02 15:04:05) | (time.Time, error) - 解析后的时间和错误 | dateutil/dateutil.go |
| dateutil | BeijingParseByFormat | 按指定格式解析时间字符串为北京时间 | format string - 时间格式, t string - 时间字符串 | (time.Time, error) - 解析后的时间和错误 | dateutil/dateutil.go |
| dateutil | GetDayAddDay | 日期加减天数计算 | lastJqEndDay string - 基准日期(格式:2006-01-02), dayNum int - 要加减的天数(可为负数) | (string, error) - 计算后的日期字符串和错误 | dateutil/dateutil.go |

## distance 模块

**文件地址**: `util/distance/distance.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| distance | CalcDistanceCheck | 计算两点间距离(带零值检查) | lat float64 - 目标纬度, lng float64 - 目标经度, orgLat float64 - 原点纬度, orgLng float64 - 原点经度 | int - 距离(米)，参数为0时返回0 | distance/distance.go |
| distance | CalcDistance | 计算两点间距离 | lat float64 - 目标纬度, lng float64 - 目标经度, orgLat float64 - 原点纬度, orgLng float64 - 原点经度 | int - 距离(米) | distance/distance.go |

## maputil 模块

**文件地址**: `util/maputil/maputil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| maputil | AddToMapMapStruct | 向双层map中添加元素 | m map[string]map[string]T - 目标双层map, key1 string - 第一层key, key2 string - 第二层key, value T - 要添加的值 | 无(原地修改) | maputil/maputil.go |

## objutil 模块

**文件地址**: `util/objutil/objutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| objutil | ExtractMethodName | 从路径中提取方法名(最后一个'/'后的内容) | path string - 完整路径字符串 | (string, error) - 提取的方法名和错误 | objutil/objutil.go |
| objutil | GetAndSetFirstReturnType | 根据方法名获取第一个返回体并实例化 | ctx context.Context - 上下文, tempService interface{} - 服务接口类型(通常为nil指针), addElem interface{} - 要设置的元素值, addElemName string - 要设置的字段名 | (interface{}, error) - 实例化后的返回体和错误 | objutil/objutil.go |

## osutil 模块

**文件地址**: `util/osutil/osutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| osutil | IsLocalEnv | 判断是否为本地环境 | 无 | bool - 是否为本地环境(DOCKER_ENV为空时为true) | osutil/osutil.go |
| osutil | GetPodIP | 获取Pod IP地址 | 无 | string - Pod的IP地址 | osutil/osutil.go |
| osutil | ListenForSIGTERM | 监听SIGTERM信号并执行回调 | ctx context.Context - 上下文, cancel context.CancelFunc - 取消函数, callback func() - 信号触发时的回调函数 | 无(启动goroutine监听) | osutil/osutil.go |
| osutil | GuessFileType | 根据文件头猜测文件类型 | buffer []byte - 文件头部字节数据(至少8字节) | (string, error) - 文件类型(pdf/jpg/png/gif/bmp/webp)和错误 | osutil/osutil.go |

## stringutil 模块

**文件地址**: `util/stringutil/stringutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| stringutil | TruncateString | 截断字符串到指定长度 | input string - 输入字符串, maxLength int - 最大长度 | string - 截断后的字符串 | stringutil/stringutil.go |
| stringutil | MD5Hash | 计算字符串的MD5哈希值 | str string - 输入字符串 | string - 32位MD5哈希值 | stringutil/stringutil.go |
| stringutil | GetFloat64Val | 字符串转float64，失败返回默认值 | str string - 要转换的字符串, defaultVal float64 - 转换失败时的默认值 | float64 - 转换结果 | stringutil/stringutil.go |
| stringutil | GetUrlParams | 解析URL参数为map | urlStr string - URL字符串 | (map[string]string, error) - 参数键值对和错误 | stringutil/stringutil.go |
| stringutil | StringToInt | 字符串转int | s string - 要转换的字符串 | (int, error) - 转换结果和错误 | stringutil/stringutil.go |
| stringutil | IntToString | 整数转字符串(支持多种整数类型) | n interface{} - 整数(支持int/int8/int16/int32/int64), defaultValue string - 转换失败时的默认值 | string - 转换后的字符串 | stringutil/stringutil.go |
| stringutil | GetPhoneMaskStr | 手机号脱敏处理(保留前3后4位) | phoneNumber string - 手机号字符串 | string - 脱敏后的手机号(如:138****1234) | stringutil/stringutil.go |
| stringutil | CheckMobile | 验证手机号格式 | phone string - 要验证的手机号 | bool - 是否为有效手机号格式 | stringutil/stringutil.go |
| stringutil | IsTimeStringYMDHMS | 验证时间字符串格式(YYYY-MM-DD HH:MM:SS) | str string - 要验证的时间字符串 | bool - 是否符合格式 | stringutil/stringutil.go |
| stringutil | IsTimeStringYMD | 验证日期字符串格式(YYYY-MM-DD) | str string - 要验证的日期字符串 | bool - 是否符合格式 | stringutil/stringutil.go |
| stringutil | IsTimeStringValid | 验证时间字符串是否符合指定格式 | timeString string - 时间字符串, format string - 时间格式 | bool - 是否符合格式 | stringutil/stringutil.go |
| stringutil | IsNil | 判断接口是否为nil或空指针 | i interface{} - 要判断的接口 | bool - 是否为nil | stringutil/stringutil.go |
| stringutil | GetJsonStrNilEmpty | 对象转JSON字符串(空值返回空字符串) | obj interface{} - 要转换的对象 | string - JSON字符串(空值返回"") | stringutil/stringutil.go |
| stringutil | GetJsonStr | 对象转JSON字符串 | obj interface{} - 要转换的对象 | string - JSON字符串 | stringutil/stringutil.go |
| stringutil | GetMaskStr | 字符串脱敏处理(支持邮箱格式) | str string - 要脱敏的字符串, plainType int - 脱敏类型(100为邮箱特殊处理) | string - 脱敏后的字符串 | stringutil/stringutil.go |
| stringutil | Substr2 | 字符串截取(支持中文) | str string - 原字符串, start int - 开始位置(从0开始), end int - 结束位置(不包含) | string - 截取后的字符串 | stringutil/stringutil.go |

## structutil 模块

**文件地址**: `util/structutil/structutil.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| structutil | GetStructArrName | 提取结构体数组中元素的类型名称 | slice interface{} - 结构体数组或切片 | string - 用逗号连接的类名列表 | structutil/structutil.go |
| structutil | GetStructName | 提取结构体类型名称 | i interface{} - 结构体实例(可为指针) | string - 类型名称 | structutil/structutil.go |
| structutil | CopyFields | 复制结构体字段值(同名字段) | src interface{} - 源结构体, dst interface{} - 目标结构体(指针) | 无(原地修改目标结构体) | structutil/structutil.go |
| structutil | CopyFieldsByParamKey | 按指定字段名复制结构体字段值，返回错误字段列表 | src interface{} - 源结构体, dst interface{} - 目标结构体(指针), reqParamKey []string - 要复制的字段名列表 | []string - 复制失败的字段名列表 | structutil/structutil.go |

## timeutil 模块

**文件地址**: `util/timeutil/time.go`

| 模块名称 | 接口名称 | 含义 | 请求参数 | 返回参数 | 文件地址 |
|---------|----------|------|----------|----------|----------|
| timeutil | TimeStr | 获取当前时间字符串(YYYY-MM-DD HH:MM:SS) | 无 | string - 当前时间字符串 | timeutil/time.go |
| timeutil | CostTms | 计算耗时(毫秒) | start time.Time - 开始时间 | int64 - 耗时毫秒数 | timeutil/time.go |

## 使用说明

### 导入方式
```go
import "util/arrayutil"
import "util/stringutil"
// 或者
import (
    "util/arrayutil"
    "util/stringutil"
    "util/dateutil"
)
```

### 注意事项
1. **泛型支持**：部分函数使用了Go 1.18+的泛型特性
2. **错误处理**：返回error的函数需要检查错误返回值
3. **时区处理**：dateutil中的函数默认使用Asia/Shanghai时区
4. **参数验证**：部分函数会对输入参数进行验证，如手机号格式、时间格式等
5. **原地修改**：部分函数会原地修改参数，如ReverseSlice、CopyFields等
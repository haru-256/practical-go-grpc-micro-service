# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common/v1/error.proto](#common_v1_error-proto)
    - [Error](#common-v1-Error)
  
- [common/v1/models.proto](#common_v1_models-proto)
    - [Category](#common-v1-Category)
    - [Product](#common-v1-Product)
  
- [command/v1/command.proto](#command_v1_command-proto)
    - [CategoryUpResult](#command-v1-CategoryUpResult)
    - [CreateCategoryRequest](#command-v1-CreateCategoryRequest)
    - [CreateCategoryResponse](#command-v1-CreateCategoryResponse)
    - [CreateProductRequest](#command-v1-CreateProductRequest)
    - [CreateProductResponse](#command-v1-CreateProductResponse)
    - [DeleteCategoryRequest](#command-v1-DeleteCategoryRequest)
    - [DeleteCategoryResponse](#command-v1-DeleteCategoryResponse)
    - [DeleteProductRequest](#command-v1-DeleteProductRequest)
    - [DeleteProductResponse](#command-v1-DeleteProductResponse)
    - [ProductUpResult](#command-v1-ProductUpResult)
    - [UpdateCategoryRequest](#command-v1-UpdateCategoryRequest)
    - [UpdateCategoryResponse](#command-v1-UpdateCategoryResponse)
    - [UpdateProductRequest](#command-v1-UpdateProductRequest)
    - [UpdateProductResponse](#command-v1-UpdateProductResponse)
  
    - [CRUD](#command-v1-CRUD)
  
    - [CategoryService](#command-v1-CategoryService)
    - [ProductService](#command-v1-ProductService)
  
- [query/v1/query.proto](#query_v1_query-proto)
    - [GetCategoryByIdRequest](#query-v1-GetCategoryByIdRequest)
    - [GetCategoryByIdResponse](#query-v1-GetCategoryByIdResponse)
    - [GetProductByIdRequest](#query-v1-GetProductByIdRequest)
    - [GetProductByIdResponse](#query-v1-GetProductByIdResponse)
    - [ListCategoriesRequest](#query-v1-ListCategoriesRequest)
    - [ListCategoriesResponse](#query-v1-ListCategoriesResponse)
    - [ListProductsRequest](#query-v1-ListProductsRequest)
    - [ListProductsResponse](#query-v1-ListProductsResponse)
    - [SearchProductsByKeywordRequest](#query-v1-SearchProductsByKeywordRequest)
    - [SearchProductsByKeywordResponse](#query-v1-SearchProductsByKeywordResponse)
    - [StreamProductsRequest](#query-v1-StreamProductsRequest)
    - [StreamProductsResponse](#query-v1-StreamProductsResponse)
  
    - [CategoryService](#query-v1-CategoryService)
    - [ProductService](#query-v1-ProductService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="common_v1_error-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/v1/error.proto



<a name="common-v1-Error"></a>

### Error



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | エラー種別 |
| message | [string](#string) |  | エラーメッセージ |





 

 

 

 



<a name="common_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common/v1/models.proto



<a name="common-v1-Category"></a>

### Category
商品カテゴリ型の定義


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | カテゴリ番号 |
| name | [string](#string) |  | カテゴリ名 |






<a name="common-v1-Product"></a>

### Product
商品型の定義


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | 商品Id |
| name | [string](#string) |  | 商品名 |
| price | [int32](#int32) |  | 単価 |
| category | [Category](#common-v1-Category) | optional | 商品カテゴリ |





 

 

 

 



<a name="command_v1_command-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/v1/command.proto



<a name="command-v1-CategoryUpResult"></a>

### CategoryUpResult
商品カテゴリ更新Result型
カテゴリ操作の結果とメタデータを返す


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [common.v1.Category](#common-v1-Category) |  | 更新されたカテゴリ情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-CreateCategoryRequest"></a>

### CreateCategoryRequest
CategoryService用のRequest/Response型


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品カテゴリ名（1-100文字） |






<a name="command-v1-CreateCategoryResponse"></a>

### CreateCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [common.v1.Category](#common-v1-Category) |  | 更新されたカテゴリ情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-CreateProductRequest"></a>

### CreateProductRequest
ProductService用のRequest/Response型


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品名（1-200文字） |
| price | [int32](#int32) |  | 単価（1円以上、999,999,999円以下） |
| category_id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |






<a name="command-v1-CreateProductResponse"></a>

### CreateProductResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | 更新された商品情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-DeleteCategoryRequest"></a>

### DeleteCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品カテゴリ名（1-100文字） |






<a name="command-v1-DeleteCategoryResponse"></a>

### DeleteCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [common.v1.Category](#common-v1-Category) |  | 更新されたカテゴリ情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-DeleteProductRequest"></a>

### DeleteProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品名（1-200文字） |
| price | [int32](#int32) |  | 単価（1円以上、999,999,999円以下） |
| category_id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |






<a name="command-v1-DeleteProductResponse"></a>

### DeleteProductResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | 更新された商品情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-ProductUpResult"></a>

### ProductUpResult
商品更新Result型
商品操作の結果とメタデータを返す


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | 更新された商品情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-UpdateCategoryRequest"></a>

### UpdateCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品カテゴリ名（1-100文字） |






<a name="command-v1-UpdateCategoryResponse"></a>

### UpdateCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [common.v1.Category](#common-v1-Category) |  | 更新されたカテゴリ情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |






<a name="command-v1-UpdateProductRequest"></a>

### UpdateProductRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crud | [CRUD](#command-v1-CRUD) |  | 更新の種類（CRUD_INSERT, CRUD_UPDATE, CRUD_DELETE） |
| id | [string](#string) |  | 商品番号（英数字、アンダースコア、ハイフンのみ） |
| name | [string](#string) |  | 商品名（1-200文字） |
| price | [int32](#int32) |  | 単価（1円以上、999,999,999円以下） |
| category_id | [string](#string) |  | 商品カテゴリ番号（英数字、アンダースコア、ハイフンのみ） |






<a name="command-v1-UpdateProductResponse"></a>

### UpdateProductResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | 更新された商品情報 |
| error | [common.v1.Error](#common-v1-Error) |  | 操作エラー情報（エラーがある場合のみ設定） |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | 操作実行時刻 |





 


<a name="command-v1-CRUD"></a>

### CRUD
更新の種類

| Name | Number | Description |
| ---- | ------ | ----------- |
| CRUD_UNSPECIFIED | 0 | 不明 |
| CRUD_INSERT | 1 | 追加 |
| CRUD_UPDATE | 2 | 変更 |
| CRUD_DELETE | 3 | 削除 |


 

 


<a name="command-v1-CategoryService"></a>

### CategoryService
商品カテゴリコマンドサービス（書き込み専用）
カテゴリのCRUD操作を提供するサービス

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateCategory | [CreateCategoryRequest](#command-v1-CreateCategoryRequest) | [CreateCategoryResponse](#command-v1-CreateCategoryResponse) | 新しい商品カテゴリを作成する |
| UpdateCategory | [UpdateCategoryRequest](#command-v1-UpdateCategoryRequest) | [UpdateCategoryResponse](#command-v1-UpdateCategoryResponse) | 既存の商品カテゴリを更新する |
| DeleteCategory | [DeleteCategoryRequest](#command-v1-DeleteCategoryRequest) | [DeleteCategoryResponse](#command-v1-DeleteCategoryResponse) | 商品カテゴリを削除する |


<a name="command-v1-ProductService"></a>

### ProductService
商品コマンドサービス型（書き込み専用）
商品のCRUD操作を提供するサービス

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateProduct | [CreateProductRequest](#command-v1-CreateProductRequest) | [CreateProductResponse](#command-v1-CreateProductResponse) | 新しい商品を作成する |
| UpdateProduct | [UpdateProductRequest](#command-v1-UpdateProductRequest) | [UpdateProductResponse](#command-v1-UpdateProductResponse) | 既存の商品を更新する |
| DeleteProduct | [DeleteProductRequest](#command-v1-DeleteProductRequest) | [DeleteProductResponse](#command-v1-DeleteProductResponse) | 商品を削除する |

 



<a name="query_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/v1/query.proto



<a name="query-v1-GetCategoryByIdRequest"></a>

### GetCategoryByIdRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | カテゴリ番号 |






<a name="query-v1-GetCategoryByIdResponse"></a>

### GetCategoryByIdResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [common.v1.Category](#common-v1-Category) |  | 商品カテゴリ |
| error | [common.v1.Error](#common-v1-Error) |  | エラー |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | タイムスタンプ |






<a name="query-v1-GetProductByIdRequest"></a>

### GetProductByIdRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | 商品番号 |






<a name="query-v1-GetProductByIdResponse"></a>

### GetProductByIdResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | 検索結果 |
| error | [common.v1.Error](#common-v1-Error) |  | 検索エラー |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | タイムスタンプ |






<a name="query-v1-ListCategoriesRequest"></a>

### ListCategoriesRequest
CategoryService用のRequest/Response型

空のリクエスト（全カテゴリ取得のため）






<a name="query-v1-ListCategoriesResponse"></a>

### ListCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| categories | [common.v1.Category](#common-v1-Category) | repeated | 商品カテゴリ複数 |
| error | [common.v1.Error](#common-v1-Error) |  | エラー |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | タイムスタンプ |






<a name="query-v1-ListProductsRequest"></a>

### ListProductsRequest
空のリクエスト（全商品取得のため）






<a name="query-v1-ListProductsResponse"></a>

### ListProductsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| products | [common.v1.Product](#common-v1-Product) | repeated | 商品複数 |
| error | [common.v1.Error](#common-v1-Error) |  | エラー |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | タイムスタンプ |






<a name="query-v1-SearchProductsByKeywordRequest"></a>

### SearchProductsByKeywordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keyword | [string](#string) |  | キーワード |






<a name="query-v1-SearchProductsByKeywordResponse"></a>

### SearchProductsByKeywordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| products | [common.v1.Product](#common-v1-Product) | repeated | 商品複数 |
| error | [common.v1.Error](#common-v1-Error) |  | エラー |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | タイムスタンプ |






<a name="query-v1-StreamProductsRequest"></a>

### StreamProductsRequest
ProductService用のRequest/Response型

空のリクエスト（全商品ストリーミング取得のため）






<a name="query-v1-StreamProductsResponse"></a>

### StreamProductsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| product | [common.v1.Product](#common-v1-Product) |  | ストリーミングされる商品 |





 

 

 


<a name="query-v1-CategoryService"></a>

### CategoryService
商品カテゴリ問合せサービス型（読み取り専用）

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCategories | [ListCategoriesRequest](#query-v1-ListCategoriesRequest) | [ListCategoriesResponse](#query-v1-ListCategoriesResponse) | すべてのカテゴリを問合せして返す |
| GetCategoryById | [GetCategoryByIdRequest](#query-v1-GetCategoryByIdRequest) | [GetCategoryByIdResponse](#query-v1-GetCategoryByIdResponse) | 指定されたIDのカテゴリを問合せして返す |


<a name="query-v1-ProductService"></a>

### ProductService
商品問合せサービス型（読み取り専用）

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StreamProducts | [StreamProductsRequest](#query-v1-StreamProductsRequest) | [StreamProductsResponse](#query-v1-StreamProductsResponse) stream | すべての商品を問合せして返す(Server streaming RPC) |
| ListProducts | [ListProductsRequest](#query-v1-ListProductsRequest) | [ListProductsResponse](#query-v1-ListProductsResponse) | すべての商品を問合せして返す |
| GetProductById | [GetProductByIdRequest](#query-v1-GetProductByIdRequest) | [GetProductByIdResponse](#query-v1-GetProductByIdResponse) | 指定されたIDの商品を問合せして返す |
| SearchProductsByKeyword | [SearchProductsByKeywordRequest](#query-v1-SearchProductsByKeywordRequest) | [SearchProductsByKeywordResponse](#query-v1-SearchProductsByKeywordResponse) | 指定されたキーワードで商品を検索して返す |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |


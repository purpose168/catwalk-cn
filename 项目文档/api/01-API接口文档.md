# Catwalk CN API 接口文档

## 概述

本文档提供了 Catwalk CN 框架的完整 API 接口文档，包括所有 HTTP 端点、请求方法、参数说明、返回格式及错误码定义。

## API 端点

### 1. 获取提供商列表（v2）

**端点**: `GET /v2/providers`

**说明**: 获取所有 AI 提供商和模型配置（推荐版本）

**请求参数**: 无

**请求头**:
| 名称 | 类型 | 说明 |
|------|------|------|
| `If-None-Match` | string | ETag 值，用于缓存 |

**响应头**:
| 名称 | 类型 | 说明 |
|------|------|------|
| `Content-Type` | string | `application/json` |
| `ETag` | string | 响应内容的 ETag |

**响应状态码**:
| 状态码 | 说明 |
|---------|------|
| `200` | 成功返回提供商列表 |
| `304` | 内容未修改（缓存） |
| `405` | 方法不允许 |
| `500` | 服务器内部错误 |

**响应示例**:

```json
[
  {
    "name": "OpenAI",
    "id": "openai",
    "type": "openai",
    "api_key": "$OPENAI_API_KEY",
    "api_endpoint": "$OPENAI_API_ENDPOINT",
    "default_large_model_id": "gpt-5.1-codex",
    "default_small_model_id": "gpt-4o",
    "models": [
      {
        "id": "gpt-5.2",
        "name": "GPT-5.2",
        "cost_per_1m_in": 1.75,
        "cost_per_1m_out": 14,
        "cost_per_1m_in_cached": 0.175,
        "cost_per_1m_out_cached": 0.175,
        "context_window": 400000,
        "default_max_tokens": 128000,
        "can_reason": true,
        "reasoning_levels": ["minimal", "low", "medium", "high"],
        "default_reasoning_effort": "medium",
        "supports_attachments": true
      }
    ]
  }
]
```

**缓存示例**:

```bash
# 首次请求
curl http://localhost:8080/v2/providers

# 获取 ETag
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"

# 后续请求使用 ETag
curl -H "If-None-Match: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
http://localhost:8080/v2/providers

# 返回 304 Not Modified
```

### 2. 获取提供商列表（已弃用）

**端点**: `GET /providers`

**说明**: 获取所有 AI 提供商和模型配置（已弃用，请使用 v2）

**请求参数**: 无

**请求头**: 无

**响应头**:
| 名称 | 类型 | 说明 |
|------|------|------|
| `Content-Type` | string | `application/json` |

**响应状态码**:
| 状态码 | 说明 |
|---------|------|
| `200` | 成功返回提供商列表 |
| `405` | 方法不允许 |
| `500` | 服务器内部错误 |

**响应示例**:

```json
[
  {
    "name": "OpenAI",
    "id": "openai",
    "type": "openai",
    "api_key": "$OPENAI_API_KEY",
    "api_endpoint": "$OPENAI_API_ENDPOINT",
    "default_large_model_id": "gpt-5.1-codex",
    "default_small_model_id": "gpt-4o",
    "models": [
      {
        "id": "gpt-5.2",
        "name": "GPT-5.2",
        "cost_per_1m_in": 1.75,
        "cost_per_1m_out": 14,
        "cost_per_1m_in_cached": 0.175,
        "cost_per_1m_out_cached": 0.175,
        "context_window": 400000,
        "default_max_tokens": 128000,
        "can_reason": true,
        "reasoning_levels": ["minimal", "low", "medium", "high"],
        "default_reasoning_effort": "medium",
        "supports_attachments": true
      }
    ]
  }
]
```

### 3. 健康检查

**端点**: `GET /healthz`

**说明**: 检查服务健康状态

**请求参数**: 无

**请求头**: 无

**响应头**: 无

**响应状态码**:
| 状态码 | 说明 |
|---------|------|
| `200` | 服务健康 |

**响应示例**:

```
OK
```

**使用示例**:

```bash
# 健康检查
curl http://localhost:8080/healthz

# 返回
OK
```

### 4. Prometheus 指标

**端点**: `GET /metrics`

**说明**: 获取 Prometheus 格式的监控指标

**请求参数**: 无

**请求头**: 无

**响应头**:
| 名称 | 类型 | 说明 |
|------|------|------|
| `Content-Type` | string | `text/plain; version=0.0.4` |

**响应状态码**:
| 状态码 | 说明 |
|---------|------|
| `200` | 成功返回指标 |

**响应示例**:

```
# HELP catwalk_providers_requests_total Total number of requests to providers endpoint
# TYPE catwalk_providers_requests_total counter
catwalk_providers_requests_total 42
```

**使用示例**:

```bash
# 获取指标
curl http://localhost:8080/metrics

# 返回
# HELP catwalk_providers_requests_total Total number of requests to providers endpoint
# TYPE catwalk_providers_requests_total counter
catwalk_providers_requests_total 42
```

## 数据模型

### Provider

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 提供商名称 |
| `id` | string | 推理提供商ID |
| `api_key` | string | API密钥（环境变量占位符） |
| `api_endpoint` | string | API端点地址（环境变量占位符） |
| `type` | string | 提供商类型 |
| `default_large_model_id` | string | 默认大模型ID |
| `default_small_model_id` | string | 默认小模型ID |
| `models` | array | 模型列表 |
| `default_headers` | object | 默认请求头 |

### Model

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 模型ID |
| `name` | string | 模型名称 |
| `cost_per_1m_in` | number | 每百万输入token的成本 |
| `cost_per_1m_out` | number | 每百万输出token的成本 |
| `cost_per_1m_in_cached` | number | 每百万缓存输入token的成本 |
| `cost_per_1m_out_cached` | number | 每百万缓存输出token的成本 |
| `context_window` | integer | 上下文窗口大小 |
| `default_max_tokens` | integer | 默认最大token数 |
| `can_reason` | boolean | 是否支持推理 |
| `reasoning_levels` | array | 推理级别列表 |
| `default_reasoning_effort` | string | 默认推理努力值 |
| `supports_attachments` | boolean | 是否支持图片附件 |
| `options` | object | 模型配置选项 |

### ModelOptions

| 字段 | 类型 | 说明 |
|------|------|------|
| `temperature` | number | 温度参数 |
| `top_p` | number | TopP参数 |
| `top_k` | integer | TopK参数 |
| `frequency_penalty` | number | 频率惩罚 |
| `presence_penalty` | number | 存在惩罚 |
| `provider_options` | object | 提供商特定选项 |

## 错误码

### HTTP 状态码

| 状态码 | 说明 | 处理建议 |
|---------|------|---------|
| `200` | 成功 | 正常处理 |
| `304` | 未修改 | 使用缓存数据 |
| `400` | 请求错误 | 检查请求参数 |
| `404` | 未找到 | 检查端点路径 |
| `405` | 方法不允许 | 检查请求方法 |
| `500` | 服务器错误 | 联系管理员 |

### 错误响应格式

```json
{
  "error": "Method not allowed"
}
```

## 使用示例

### cURL

```bash
# 获取提供商列表
curl http://localhost:8080/v2/providers

# 使用 ETag 缓存
curl -H "If-None-Match: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
http://localhost:8080/v2/providers

# 健康检查
curl http://localhost:8080/healthz

# 获取指标
curl http://localhost:8080/metrics
```

### Go

```go
import (
    "encoding/json"
    "net/http"
)

type Provider struct {
    Name   string `json:"name"`
    ID     string `json:"id"`
    Models []Model `json:"models"`
}

type Model struct {
    ID   string  `json:"id"`
    Name string  `json:"name"`
}

func getProviders() ([]Provider, error) {
    resp, err := http.Get("http://localhost:8080/v2/providers")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var providers []Provider
    if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
        return nil, err
    }

    return providers, nil
}
```

### JavaScript

```javascript
async function getProviders() {
    const response = await fetch('http://localhost:8080/v2/providers');
    const providers = await response.json();
    return providers;
}

async function getProvidersWithETag(etag) {
    const response = await fetch('http://localhost:8080/v2/providers', {
        headers: {
            'If-None-Match': etag
        }
    });
    
    if (response.status === 304) {
        return null; // 使用缓存
    }
    
    return await response.json();
}
```

### Python

```python
import requests

def get_providers():
    response = requests.get('http://localhost:8080/v2/providers')
    return response.json()

def get_providers_with_etag(etag):
    headers = {'If-None-Match': etag}
    response = requests.get('http://localhost:8080/v2/providers', headers=headers)
    
    if response.status_code == 304:
        return None  # 使用缓存
    
    return response.json()
```

## 性能优化

### 使用 ETag 缓存

```bash
# 首次请求
curl -I http://localhost:8080/v2/providers

# 获取 ETag
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"

# 后续请求使用 ETag
curl -H "If-None-Match: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
http://localhost:8080/v2/providers

# 返回 304 Not Modified，减少网络传输
```

### 使用 HEAD 请求

```bash
# 只检查 ETag
curl -I http://localhost:8080/v2/providers

# 返回头信息
HTTP/1.1 200 OK
Content-Type: application/json
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

## 总结

Catwalk CN 提供简洁而强大的 API：

1. **RESTful 设计**: 遵循 REST 原则
2. **JSON 格式**: 使用标准 JSON 格式
3. **ETag 缓存**: 支持客户端缓存
4. **健康检查**: 提供健康检查端点
5. **Prometheus 集成**: 支持监控指标
6. **版本控制**: 支持 API 版本

这种 API 设计使得 Catwalk CN 成为一个易于集成和使用的 AI 提供商数据库服务。

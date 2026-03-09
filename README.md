# Omniwriter

> ⚠️ **Disclaimer:** This project was largely created via AI-assisted development.
> It is an experimental project and may contain security quirks or require further optimization.

**Multi-format data transformation library for Go**

Omniwriter extends [omniparser](https://github.com/jf-tech/omniparser) from a JSON-output parser into a comprehensive transformation library supporting conversions between JSON, CSV, EDI, XML, and text formats.

## ✨ Features

- **Full Format Support**: Transform between JSON, CSV, EDI (X12/EDIFACT), XML, and text
- **Schema-Driven**: JSON schemas control all transformations using familiar omniparser patterns
- **Passthrough Optimization**: Same-format transforms use fast path

## 🚀 Quick Start

```go
import "github.com/pbagtoltol/omniwriter/pkg/omniwriter"

schema, _ := os.ReadFile("transform.json")
input, _ := os.Open("data.edi")

result, err := omniwriter.Transform(context.Background(),
    omniwriter.TransformRequest{
        SourceFormat: omniwriter.FormatEDI,
        TargetFormat: omniwriter.FormatCSV,
        Mapping:      schema,
        Input:        input,
    })

if err != nil {
    log.Fatal(err)
}

os.WriteFile("output.csv", result.Output, 0644)
```

## 📊 Supported Transformations

```
          JSON   CSV  EDI  XML  Text
JSON       ✅    ✅   ✅   ✅    ✅
CSV        ✅    ✅   ✅   ✅    ✅
EDI        ✅    ✅   ✅   ✅    ✅
XML        ✅    ✅   ✅   ✅    ✅
```

## 📚 Documentation

- [Implementation Plan](docs/IMPLEMENTATION_PLAN.md) - Architecture and roadmap
- [Project Summary](docs/PROJECT_SUMMARY.md) - Current status and capabilities
- [Implementation Checklist](docs/IMPLEMENTATION_CHECKLIST.md) - Detailed progress tracking
- [Examples](examples/) - 8 working examples with READMEs

## 🎯 Use Cases

- **EDI Processing**: Parse and convert between JSON, CSV, EDI (X12/EDIFACT), XML, and text
- **Data Integration**: Transform between different system formats
- **ETL Pipelines**: Schema-driven data transformations

## 🏗️ Status

- ✅ Core engine and API
- ✅ All format emitters (CSV, JSON, EDI, XML, text)
- ✅ Comprehensive test coverage (42 tests)
- ✅ Working examples
- ⏳ Performance benchmarks (optional)
- ⏳ Custom format plugins (optional)

## 📄 License

MIT

## 🙏 Acknowledgments

Built on [jf-tech/omniparser](https://github.com/jf-tech/omniparser) - excellent parser foundation for schema-driven data processing.

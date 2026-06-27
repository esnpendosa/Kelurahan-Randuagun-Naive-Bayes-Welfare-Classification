# Design Document: Dynamic Training Model Enhancement

## Overview

This design document outlines the technical architecture for transforming the existing web-based Naive Bayes welfare classification system into a native Windows desktop application with dynamic training model capabilities, enhanced filtering, and automatic data population features.

## Critical Design Priorities

1. **Calculation Accuracy**: Maintain 100% computational equivalence with Excel reference
2. **Desktop Native**: Migrate from web browser to Wails v2 desktop application
3. **Feature Completeness**: Implement all requested enhancements

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Windows Desktop Application                 │
│                    (Wails v2)                           │
├─────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────┐ │
│  │         Frontend (HTML/CSS/JavaScript)             │ │
│  │  - UI Components (existing pages)                  │ │
│  │  - Dynamic Model Selection                         │ │
│  │  - Filter Controls                                 │ │
│  │  - Data Assignment Interface                       │ │
│  └───────────────────────────────────────────────────┘ │
│                          ↕ (Wails Bindings)             │
│  ┌───────────────────────────────────────────────────┐ │
│  │           Backend (Go)                             │ │
│  │  ┌─────────────────────────────────────────────┐  │ │
│  │  │  Application Layer (main.go)                │  │ │
│  │  │  - Wails Runtime                            │  │ │
│  │  │  - Route Handlers → Service Methods         │  │ │
│  │  └─────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────┐  │ │
│  │  │  Service Layer                              │  │ │
│  │  │  - ModelService (dual splits)               │  │ │
│  │  │  - ClassificationService                    │  │ │
│  │  │  - EvaluationService                        │  │ │
│  │  │  - WargaService                             │  │ │
│  │  └─────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────┐  │ │
│  │  │  Classifier Package                         │  │ │
│  │  │  - NaiveBayes (enhanced)                    │  │ │
│  │  │  - Evaluation Metrics                       │  │ │
│  │  └─────────────────────────────────────────────┘  │ │
│  │  ┌─────────────────────────────────────────────┐  │ │
│  │  │  Database Layer (internal/db)               │  │ │
│  │  │  - SQLite Connection                        │  │ │
│  │  │  - Enhanced Schema                          │  │ │
│  │  └─────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          ↕
              ┌─────────────────────┐
              │  SQLite Database    │
              │  (data_skripsi.db)  │
              └─────────────────────┘
```

## Database Design

### Schema Changes

#### 1. Enhanced `warga` Table

```sql
ALTER TABLE warga ADD COLUMN data_latih_2 INTEGER DEFAULT 0;
```

**Fields:**
- `id`: INTEGER PRIMARY KEY
- `nik`: TEXT (existing)
- `nama`: TEXT (existing)
- `alamat`: TEXT (existing)
- `label_kelas`: INTEGER (existing - 1 to 6)
- `data_latih`: INTEGER (existing - 0=test, 1=training for Split 1)
- `data_latih_2`: INTEGER (NEW - 0=test, 1=training for Split 2)

#### 2. Existing Tables (no changes)

- `data_indikator`: Stores 36 indicators (IM1-IM36) per citizen
- `hasil_klasifikasi`: Stores classification results with probability distribution
- `users`: Authentication (simplified to single logout)

### Migration Strategy

```go
func MigrateToV2(db *sql.DB) error {
    // Check if data_latih_2 exists
    var count int
    err := db.QueryRow(`
        SELECT COUNT(*) FROM pragma_table_info('warga') 
        WHERE name='data_latih_2'
    `).Scan(&count)
    
    if count == 0 {
        // Add column if missing
        _, err = db.Exec("ALTER TABLE warga ADD COLUMN data_latih_2 INTEGER DEFAULT 0")
        if err != nil {
            return err
        }
    }
    return nil
}
```

## Backend Design

### Component Structure

#### 1. Desktop Application Entry (main.go)

**Current State:** HTTP server with Gin router
**New State:** Wails v2 application

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // Initialize database
    database := db.InitDB("data_skripsi.db")
    defer database.Close()
    
    // Run migration
    db.MigrateToV2(database)
    
    // Create application services
    app := &App{
        db: database,
        modelService: services.NewModelService(database),
        classificationService: services.NewClassificationService(database),
        evaluationService: services.NewEvaluationService(database),
        wargaService: services.NewWargaService(database),
    }
    
    // Create Wails application
    err := wails.Run(&options.App{
        Title:  "Sistem Klasifikasi Kesejahteraan Randuagung",
        Width:  1280,
        Height: 800,
        MinWidth: 1024,
        MinHeight: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
        OnStartup:  app.startup,
        OnShutdown: app.shutdown,
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
            WindowIsTranslucent:  false,
            DisableWindowIcon:    false,
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

#### 2. Service Layer Architecture

**services/model_service.go** - Dual Training Model Management

```go
type ModelService struct {
    db *sql.DB
    model1 *classifier.KlasifikasiNaiveBayes // Training Split 1
    model2 *classifier.KlasifikasiNaiveBayes // Training Split 2
    activeModel int // 1 or 2
}

// LoadTrainingData1 loads data where data_latih = 1
func (s *ModelService) LoadTrainingData1() ([]map[string]string, []int, error)

// LoadTestData1 loads data where data_latih = 0
func (s *ModelService) LoadTestData1() ([]map[string]string, []int, error)

// LoadTrainingData2 loads data where data_latih_2 = 1
func (s *ModelService) LoadTrainingData2() ([]map[string]string, []int, error)

// LoadTestData2 loads data where data_latih_2 = 0
func (s *ModelService) LoadTestData2() ([]map[string]string, []int, error)

// TrainModel trains the selected model (1 or 2)
func (s *ModelService) TrainModel(modelNum int) error

// GetDataDistribution returns training/test counts for both splits
func (s *ModelService) GetDataDistribution() (DataDistribution, error)

// SetActiveModel selects which model to use for classification
func (s *ModelService) SetActiveModel(modelNum int) error

// GetActiveModel returns currently active model
func (s *ModelService) GetActiveModel() *classifier.KlasifikasiNaiveBayes
```

**services/evaluation_service.go** - Performance Metrics

```go
type EvaluationService struct {
    db *sql.DB
}

type EvaluationMetrics struct {
    Accuracy      float64
    ConfusionMatrix [6][6]int
    Precision     [6]float64
    Recall        [6]float64
    F1Score       [6]float64
    MacroPrecision float64
    MacroRecall    float64
    MacroF1        float64
}

// EvaluateModel runs test data through model and calculates metrics
func (s *EvaluationService) EvaluateModel(
    model *classifier.KlasifikasiNaiveBayes,
    testData []map[string]string,
    testLabels []int,
) (EvaluationMetrics, error)

// FilterByClass filters evaluation results for specific welfare class
func (s *EvaluationService) FilterByClass(
    metrics EvaluationMetrics,
    class int,
) (EvaluationMetrics, error)
```

**services/classification_service.go** - Citizen Classification

```go
type ClassificationService struct {
    db *sql.DB
    modelService *ModelService
}

// GetCitizenIndicators retrieves existing indicator values
func (s *ClassificationService) GetCitizenIndicators(wargaID int) (map[string]string, error)

// ClassifyCitizen performs classification using active model
func (s *ClassificationService) ClassifyCitizen(
    wargaID int,
    indicators map[string]string,
) (ClassificationResult, error)

// SaveClassificationResult stores prediction in database
func (s *ClassificationService) SaveClassificationResult(
    wargaID int,
    result ClassificationResult,
) error
```

**services/warga_service.go** - Citizen Data Management

```go
type WargaService struct {
    db *sql.DB
}

// GetAllWarga retrieves all citizens with optional filter
func (s *WargaService) GetAllWarga(filter WargaFilter) ([]Warga, error)

// UpdateDataRole updates data_latih_2 value for a citizen
func (s *WargaService) UpdateDataRole(wargaID int, role int) error

// GetWargaByClass filters citizens by welfare category
func (s *WargaService) GetWargaByClass(class int) ([]Warga, error)
```

#### 3. Enhanced Classifier Package

**internal/classifier/naive_bayes.go** - Enhanced Model

```go
// Add verification method for Excel comparison
func (nb *KlasifikasiNaiveBayes) VerifyAgainstExcel(
    testCase map[string]string,
    expectedProbs map[int]float64,
    tolerance float64,
) (bool, error) {
    predicted := nb.Prediksi(testCase)
    
    for class, expectedProb := range expectedProbs {
        actualProb := predicted[KelasKesejahteraan(class)]
        diff := math.Abs(actualProb - expectedProb)
        
        if diff > tolerance {
            return false, fmt.Errorf(
                "Class %d: expected %.6f, got %.6f (diff: %.6f)",
                class, expectedProb, actualProb, diff,
            )
        }
    }
    return true, nil
}

// Export model state for debugging
func (nb *KlasifikasiNaiveBayes) ExportState() ModelState {
    return ModelState{
        Prior: nb.PeluangPrior,
        Likelihood: nb.PeluangLikelihood,
        Features: nb.DaftarFitur,
    }
}
```

**internal/classifier/evaluation.go** - New File

```go
package classifier

type ConfusionMatrix [6][6]int

func CalculateAccuracy(cm ConfusionMatrix) float64 {
    correct := 0
    total := 0
    for i := 0; i < 6; i++ {
        for j := 0; j < 6; j++ {
            total += cm[i][j]
            if i == j {
                correct += cm[i][j]
            }
        }
    }
    return float64(correct) / float64(total)
}

func CalculatePrecision(cm ConfusionMatrix, class int) float64 {
    truePositive := cm[class][class]
    predictedPositive := 0
    for i := 0; i < 6; i++ {
        predictedPositive += cm[i][class]
    }
    if predictedPositive == 0 {
        return 0
    }
    return float64(truePositive) / float64(predictedPositive)
}

func CalculateRecall(cm ConfusionMatrix, class int) float64 {
    truePositive := cm[class][class]
    actualPositive := 0
    for j := 0; j < 6; j++ {
        actualPositive += cm[class][j]
    }
    if actualPositive == 0 {
        return 0
    }
    return float64(truePositive) / float64(actualPositive)
}

func CalculateF1Score(precision, recall float64) float64 {
    if precision+recall == 0 {
        return 0
    }
    return 2 * (precision * recall) / (precision + recall)
}
```

## Frontend Design

### Component Updates

#### 1. Training Model Page Enhancement

**File:** `frontend/training-model.html`

**New Features:**
- Tabbed interface for Split 1 and Split 2
- Data distribution charts
- Model selection dropdown
- Process calculation button
- Evaluation results display

**Layout:**
```html
<div class="training-container">
    <!-- Tab Navigation -->
    <ul class="nav nav-tabs">
        <li class="active"><a data-toggle="tab" href="#split1">Data Split 1</a></li>
        <li><a data-toggle="tab" href="#split2">Data Split 2</a></li>
    </ul>
    
    <!-- Tab Content -->
    <div class="tab-content">
        <!-- Split 1 Panel -->
        <div id="split1" class="tab-pane active">
            <div class="data-distribution">
                <h4>Distribusi Data Training 1</h4>
                <canvas id="chart-training-1"></canvas>
                <p>Total: <span id="count-training-1"></span></p>
            </div>
            <div class="data-distribution">
                <h4>Distribusi Data Uji 1</h4>
                <canvas id="chart-test-1"></canvas>
                <p>Total: <span id="count-test-1"></span></p>
            </div>
        </div>
        
        <!-- Split 2 Panel -->
        <div id="split2" class="tab-pane">
            <div class="data-distribution">
                <h4>Distribusi Data Training 2</h4>
                <canvas id="chart-training-2"></canvas>
                <p>Total: <span id="count-training-2"></span></p>
            </div>
            <div class="data-distribution">
                <h4>Distribusi Data Uji 2</h4>
                <canvas id="chart-test-2"></canvas>
                <p>Total: <span id="count-test-2"></span></p>
            </div>
            <div class="dynamic-assignment">
                <button id="btn-edit-roles">Edit Data Roles</button>
            </div>
        </div>
    </div>
    
    <!-- Model Selection and Processing -->
    <div class="model-selection">
        <label>Pilih Model:</label>
        <select id="select-model">
            <option value="1">Training Model 1</option>
            <option value="2">Training Model 2</option>
        </select>
        <button id="btn-process" class="btn btn-primary">Proses Hitung</button>
    </div>
    
    <!-- Evaluation Results -->
    <div id="evaluation-results" style="display:none;">
        <h3>Hasil Evaluasi</h3>
        <div class="metrics">
            <p>Accuracy: <span id="accuracy"></span></p>
            <table id="confusion-matrix"></table>
            <table id="class-metrics"></table>
        </div>
    </div>
</div>
```

#### 2. Classification Page Enhancement

**File:** `frontend/klasifikasi.html`

**Enhancement:** Auto-populate indicators

```javascript
// When citizen is selected
async function onCitizenSelect(wargaID) {
    // Call backend to get existing indicators
    const indicators = await window.go.main.App.GetCitizenIndicators(wargaID);
    
    if (indicators) {
        // Auto-fill form fields
        for (let field in indicators) {
            document.getElementById(field).value = indicators[field];
            // Add visual indication
            document.getElementById(field).classList.add('auto-filled');
        }
        
        // Enable submit button
        document.getElementById('btn-classify').disabled = false;
    }
}
```

#### 3. Reports Page Enhancement

**File:** `frontend/reports.html`

**New Feature:** Welfare category filter

```html
<div class="filter-controls">
    <label>Filter Kategori:</label>
    <select id="filter-class">
        <option value="">Semua Kategori</option>
        <option value="1">Sangat Miskin</option>
        <option value="2">Miskin</option>
        <option value="3">Hampir Miskin</option>
        <option value="4">Rentan Miskin</option>
        <option value="5">Pas-pasan</option>
        <option value="6">Menengah ke Atas</option>
    </select>
    <button id="btn-apply-filter">Terapkan</button>
    <button id="btn-clear-filter">Hapus Filter</button>
</div>

<div id="filtered-results">
    <table id="warga-table"></table>
    <div id="statistics"></div>
</div>
```

#### 4. Warga Data Page Simplification

**File:** `frontend/warga.html`

**Change:** Remove training/test labels from display

```html
<!-- OLD: -->
<td>
    {{if .DataLatih}}Data Latih{{else}}Data Uji{{end}}
</td>

<!-- NEW: Remove this column entirely -->
```

#### 5. Dynamic Role Assignment Modal

**New Component:** `frontend/components/role-assignment-modal.html`

```html
<div id="role-assignment-modal" class="modal">
    <div class="modal-content">
        <h3>Atur Role Data untuk Split 2</h3>
        <table id="role-table">
            <thead>
                <tr>
                    <th>NIK</th>
                    <th>Nama</th>
                    <th>Kelas</th>
                    <th>Role Split 2</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody id="role-tbody">
                <!-- Dynamically populated -->
            </tbody>
        </table>
        <button id="btn-save-roles">Simpan</button>
        <button id="btn-cancel">Batal</button>
    </div>
</div>
```

### Wails Bindings (JavaScript ↔ Go)

```javascript
// Frontend calls backend methods via Wails runtime

// Get data distribution
const distribution = await window.go.main.App.GetDataDistribution();

// Train model
await window.go.main.App.TrainModel(modelNum);

// Evaluate model
const metrics = await window.go.main.App.EvaluateModel(modelNum);

// Get citizen indicators
const indicators = await window.go.main.App.GetCitizenIndicators(wargaID);

// Classify citizen
const result = await window.go.main.App.ClassifyCitizen(wargaID, indicators);

// Update data role
await window.go.main.App.UpdateDataRole(wargaID, role);

// Get filtered warga
const warga = await window.go.main.App.GetWargaByClass(classNum);
```

## Build and Deployment

### Development Setup

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Initialize Wails in existing project
wails init -n Klasifikasi-Warga-Randuagung

# Run in development mode with hot reload
wails dev
```

### Production Build

```bash
# Build Windows executable
wails build -platform windows/amd64

# Output: build/bin/Klasifikasi-Warga-Randuagung.exe
```

### Distribution Package

```
Klasifikasi-Warga-Randuagung/
├── Klasifikasi-Warga-Randuagung.exe
├── data_skripsi.db (optional, auto-created)
├── README.txt
└── icon.ico
```

## Testing Strategy

### 1. Calculation Verification Tests

```go
func TestExcelEquivalence(t *testing.T) {
    // Load test cases from Excel
    testCases := loadExcelTestCases("klasifikasi naive bayes tambahan data.xlsx")
    
    // Train model with Split 1
    model := trainModelFromExcel("Data Training 1")
    
    // Test each case from "Data Uji 1"
    for _, tc := range testCases {
        predicted := model.Prediksi(tc.Indicators)
        
        // Verify class prediction
        assert.Equal(t, tc.ExpectedClass, model.AmbilKelasTerbaik(predicted))
        
        // Verify probabilities within tolerance
        for class, expectedProb := range tc.ExpectedProbs {
            actualProb := predicted[class]
            assert.InDelta(t, expectedProb, actualProb, 0.00001)
        }
    }
}
```

### 2. Integration Tests

- Database migration
- Model training/evaluation cycle
- Frontend-backend communication via Wails
- Filter operations
- Role assignment updates

### 3. UI Tests

- Modal interactions
- Tab switching
- Form auto-population
- Filter application

## Migration Checklist

- [ ] Install Wails v2
- [ ] Add Wails configuration
- [ ] Refactor main.go for Wails runtime
- [ ] Add data_latih_2 migration
- [ ] Implement ModelService with dual splits
- [ ] Implement EvaluationService
- [ ] Implement ClassificationService enhancements
- [ ] Implement WargaService filters
- [ ] Update frontend with Wails bindings
- [ ] Add tabbed training interface
- [ ] Add model selection UI
- [ ] Add filter controls
- [ ] Add role assignment modal
- [ ] Remove training labels from warga display
- [ ] Consolidate logout button
- [ ] Test Excel calculation equivalence
- [ ] Build production executable
- [ ] User acceptance testing

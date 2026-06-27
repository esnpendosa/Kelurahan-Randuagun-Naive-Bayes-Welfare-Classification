# Implementation Tasks: Dynamic Training Model Enhancement

## Phase 1: Desktop Application Migration (Wails v2)

### Task 1.1: Setup Wails Development Environment
**Priority:** High | **Estimate:** 2 hours

- [ ] Install Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [ ] Install WebView2 runtime (if not present on Windows)
- [ ] Verify Go version compatibility (1.18+)
- [ ] Test Wails installation with hello-world project

**Acceptance:**
- Wails CLI runs successfully
- Can create and run a basic Wails application

---

### Task 1.2: Initialize Wails in Existing Project
**Priority:** High | **Estimate:** 3 hours

- [ ] Create `wails.json` configuration file
- [ ] Restructure project folders for Wails convention:
  - Move HTML/CSS/JS to `frontend/` directory
  - Keep Go code in root/internal
- [ ] Create `frontend/index.html` as main entry point
- [ ] Configure embed directives for assets

**Acceptance:**
- Project structure matches Wails conventions
- `wails dev` command launches successfully
- Frontend assets load in webview window

---

### Task 1.3: Refactor main.go for Wails Runtime
**Priority:** High | **Estimate:** 4 hours

- [ ] Remove Gin/HTTP server code
- [ ] Create App struct with service dependencies
- [ ] Implement startup() and shutdown() methods
- [ ] Configure Wails.Run() with window options
- [ ] Set application title and icon
- [ ] Configure window size (1280x800, min 1024x768)

**Acceptance:**
- Application launches as native window (not browser)
- Window has custom title and icon
- Window controls (minimize, maximize, close) work
- Database initializes on startup

---

### Task 1.4: Create Wails Bindings for Existing Routes
**Priority:** High | **Estimate:** 5 hours

Convert HTTP handlers to Go methods callable from JavaScript:

- [ ] Authentication methods (Login, Logout)
- [ ] Warga CRUD methods
- [ ] Indicator data methods
- [ ] Classification methods
- [ ] Report generation methods

**Acceptance:**
- All backend functions accessible via `window.go.main.App.*`
- JavaScript can call Go methods successfully
- Return types properly serialized to JSON

---

### Task 1.5: Update Frontend JavaScript for Wails
**Priority:** High | **Estimate:** 4 hours

- [ ] Replace axios/fetch calls with Wails runtime calls
- [ ] Update all API endpoints to use `window.go.*` pattern
- [ ] Test form submissions
- [ ] Test data retrieval and display
- [ ] Handle errors from Go methods

**Acceptance:**
- All existing pages work in desktop window
- Form submissions succeed
- Data loads and displays correctly
- Error handling works

---

## Phase 2: Database Schema Enhancement

### Task 2.1: Create Database Migration
**Priority:** High | **Estimate:** 2 hours

- [ ] Create `internal/db/migrations.go`
- [ ] Implement `MigrateToV2()` function
- [ ] Add `data_latih_2` column to warga table
- [ ] Test migration on existing database
- [ ] Test migration on fresh database

**Acceptance:**
- Migration adds `data_latih_2` column with default 0
- Migration is idempotent (safe to run multiple times)
- Existing data preserved
- No errors on fresh or existing databases

---

### Task 2.2: Update Database Access Layer
**Priority:** Medium | **Estimate:** 2 hours

- [ ] Update warga queries to include `data_latih_2`
- [ ] Add query for filtering by `data_latih_2`
- [ ] Update insert/update statements for warga table
- [ ] Add indexes if needed for performance

**Acceptance:**
- Can query citizens by `data_latih_2` value
- Can update `data_latih_2` for individual citizens
- Queries perform efficiently

---

## Phase 3: Backend Service Implementation

### Task 3.1: Implement ModelService (Dual Training Splits)
**Priority:** High | **Estimate:** 6 hours

Create `services/model_service.go`:

- [ ] Create ModelService struct with two model instances
- [ ] Implement `LoadTrainingData1()` - query where data_latih=1
- [ ] Implement `LoadTestData1()` - query where data_latih=0
- [ ] Implement `LoadTrainingData2()` - query where data_latih_2=1
- [ ] Implement `LoadTestData2()` - query where data_latih_2=0
- [ ] Implement `TrainModel(modelNum int)` - trains model 1 or 2
- [ ] Implement `SetActiveModel(modelNum int)` - selects active model
- [ ] Implement `GetActiveModel()` - returns current model
- [ ] Implement `GetDataDistribution()` - counts for all splits

**Acceptance:**
- Can load data for both splits independently
- Can train model 1 or model 2
- Can select which model is active
- Data distribution returns accurate counts

---

### Task 3.2: Implement EvaluationService
**Priority:** High | **Estimate:** 5 hours

Create `services/evaluation_service.go`:

- [ ] Create EvaluationService struct
- [ ] Implement `EvaluateModel()` method
- [ ] Build confusion matrix from predictions
- [ ] Calculate accuracy
- [ ] Calculate precision per class
- [ ] Calculate recall per class
- [ ] Calculate F1-score per class
- [ ] Calculate macro-averaged metrics
- [ ] Implement `FilterByClass()` for single-class evaluation

**Acceptance:**
- Confusion matrix correctly tallies predictions
- Metrics match manual calculations
- Macro averages correct
- Class filtering works

---

### Task 3.3: Implement Enhanced ClassificationService
**Priority:** High | **Estimate:** 4 hours

Create `services/classification_service.go`:

- [ ] Create ClassificationService struct
- [ ] Implement `GetCitizenIndicators(wargaID)` - retrieves existing data
- [ ] Implement `ClassifyCitizen()` - uses active model
- [ ] Implement `SaveClassificationResult()` - stores to database
- [ ] Integrate with ModelService for active model selection

**Acceptance:**
- Can retrieve indicator values for existing citizens
- Returns empty map for citizens without indicators
- Classification uses the currently active model
- Results saved to database correctly

---

### Task 3.4: Implement WargaService with Filtering
**Priority:** Medium | **Estimate:** 3 hours

Create `services/warga_service.go`:

- [ ] Create WargaService struct
- [ ] Implement `GetAllWarga(filter)` with optional class filter
- [ ] Implement `GetWargaByClass(class)` - filters by label_kelas
- [ ] Implement `UpdateDataRole(wargaID, role)` - updates data_latih_2
- [ ] Implement `GetWargaWithIndicators()` - joins with data_indikator

**Acceptance:**
- Can retrieve all citizens
- Can filter by welfare class
- Can update data_latih_2 values
- Filtering performs efficiently

---

### Task 3.5: Add Excel Verification Method to Classifier
**Priority:** High | **Estimate:** 3 hours

Update `internal/classifier/naive_bayes.go`:

- [ ] Add `VerifyAgainstExcel()` method
- [ ] Compare predicted probabilities with expected values
- [ ] Return detailed comparison report
- [ ] Add tolerance parameter (default 0.00001)
- [ ] Add `ExportState()` for debugging

**Acceptance:**
- Can compare predictions against Excel reference
- Reports differences accurately
- Tolerance parameter works
- Useful for testing

---

### Task 3.6: Create Evaluation Metrics Package
**Priority:** Medium | **Estimate:** 3 hours

Create `internal/classifier/evaluation.go`:

- [ ] Implement `CalculateAccuracy()`
- [ ] Implement `CalculatePrecision()`
- [ ] Implement `CalculateRecall()`
- [ ] Implement `CalculateF1Score()`
- [ ] Add helper functions for confusion matrix operations

**Acceptance:**
- All metric functions return correct values
- Functions tested against known examples
- Confusion matrix operations work correctly

---

## Phase 4: Frontend Enhancement

### Task 4.1: Create Tabbed Training Model Interface
**Priority:** High | **Estimate:** 4 hours

Update `frontend/training-model.html`:

- [ ] Add Bootstrap tabs for Split 1 and Split 2
- [ ] Create data distribution display sections
- [ ] Add Chart.js for visualizing distributions
- [ ] Display training and test counts
- [ ] Add model selection dropdown
- [ ] Add "Proses Hitung" button
- [ ] Wire up Wails bindings for data loading

**Acceptance:**
- Tabs switch between Split 1 and Split 2
- Charts display class distributions
- Counts update correctly
- Model selection persists

---

### Task 4.2: Implement Model Selection and Evaluation UI
**Priority:** High | **Estimate:** 5 hours

Continue `frontend/training-model.html`:

- [ ] Wire "Proses Hitung" button to backend
- [ ] Show loading indicator during training
- [ ] Display evaluation results section
- [ ] Render confusion matrix as table
- [ ] Display accuracy, precision, recall, F1 metrics
- [ ] Format metrics as percentages
- [ ] Add export results button (optional)

**Acceptance:**
- Button triggers model training
- Loading indicator shows during processing
- Results display after completion
- All metrics visible and formatted
- Confusion matrix readable

---

### Task 4.3: Create Dynamic Role Assignment Modal
**Priority:** Medium | **Estimate:** 5 hours

Create `frontend/components/role-assignment-modal.html`:

- [ ] Create modal overlay and content
- [ ] Fetch all citizens with current roles
- [ ] Display in sortable/filterable table
- [ ] Add toggle switches for training/test role
- [ ] Implement save functionality
- [ ] Show success/error messages
- [ ] Wire to WargaService.UpdateDataRole()

**Acceptance:**
- Modal opens from Split 2 tab
- Table shows all citizens with current Split 2 roles
- Can toggle individual citizen roles
- Save button updates database
- Changes reflect immediately in distribution charts

---

### Task 4.4: Implement Auto-fill Classification Form
**Priority:** High | **Estimate:** 3 hours

Update `frontend/klasifikasi.html`:

- [ ] Add citizen selection dropdown
- [ ] Wire onChange to GetCitizenIndicators()
- [ ] Auto-populate form fields when data exists
- [ ] Add visual styling for auto-filled fields
- [ ] Enable submit button when form complete
- [ ] Allow manual editing of auto-filled fields

**Acceptance:**
- Selecting citizen loads indicators if they exist
- Fields populate automatically
- Auto-filled fields visually distinct
- Can still manually edit fields
- Submit button enables appropriately

---

### Task 4.5: Add Welfare Category Filter to Reports
**Priority:** Medium | **Estimate:** 3 hours

Update `frontend/reports.html`:

- [ ] Add filter dropdown with 6 welfare categories
- [ ] Add "Apply Filter" button
- [ ] Add "Clear Filter" button
- [ ] Wire filter to WargaService.GetWargaByClass()
- [ ] Update statistics for filtered data
- [ ] Show active filter indicator

**Acceptance:**
- Filter dropdown shows all 6 categories
- Apply button filters table and statistics
- Clear button resets to all data
- Active filter clearly indicated
- Statistics recalculate for filtered subset

---

### Task 4.6: Remove Training Labels from Warga Display
**Priority:** Low | **Estimate:** 1 hour

Update `frontend/warga.html`:

- [ ] Remove "Data Latih/Uji" column from table
- [ ] Update table headers
- [ ] Test that table still displays correctly

**Acceptance:**
- Table no longer shows training/test labels
- All other columns still visible
- Styling consistent

---

### Task 4.7: Consolidate Logout Button
**Priority:** Low | **Estimate:** 1 hour

Update navigation:

- [ ] Remove duplicate logout/exit buttons
- [ ] Keep single "Logout" button
- [ ] Test logout functionality
- [ ] Ensure session clears properly

**Acceptance:**
- Only one logout button visible
- Logout works correctly
- Redirects to login page

---

## Phase 5: Testing and Verification

### Task 5.1: Excel Calculation Verification Tests
**Priority:** Critical | **Estimate:** 6 hours

Create `internal/classifier/naive_bayes_test.go`:

- [ ] Load "klasifikasi naive bayes tambahan data.xlsx"
- [ ] Extract training data from "Data Training 1" sheet
- [ ] Extract test data from "Data Uji 1" sheet
- [ ] Train model with extracted data
- [ ] Compare predicted classes with Excel
- [ ] Compare probability distributions (tolerance 0.001%)
- [ ] Test with Training/Test Split 2 as well
- [ ] Document any discrepancies

**Acceptance:**
- Tests pass with 100% accuracy match to Excel
- Probabilities within 0.001% tolerance
- Confusion matrix matches Excel
- All 6 classes tested

---

### Task 5.2: Model Service Integration Tests
**Priority:** High | **Estimate:** 4 hours

Create `services/model_service_test.go`:

- [ ] Test LoadTrainingData1() returns correct count
- [ ] Test LoadTestData1() returns correct count
- [ ] Test LoadTrainingData2() returns correct count
- [ ] Test LoadTestData2() returns correct count
- [ ] Test TrainModel() completes without errors
- [ ] Test SetActiveModel() switches models
- [ ] Test GetDataDistribution() returns accurate counts

**Acceptance:**
- All data loading tests pass
- Model training succeeds
- Model switching works
- Distribution counts accurate

---

### Task 5.3: Evaluation Service Tests
**Priority:** High | **Estimate:** 3 hours

Create `services/evaluation_service_test.go`:

- [ ] Test confusion matrix calculation with known data
- [ ] Test accuracy calculation
- [ ] Test precision calculation per class
- [ ] Test recall calculation per class
- [ ] Test F1-score calculation
- [ ] Test macro-averaging
- [ ] Test FilterByClass()

**Acceptance:**
- All metrics calculate correctly
- Results match manual calculations
- Filtering works

---

### Task 5.4: Frontend-Backend Integration Tests
**Priority:** High | **Estimate:** 4 hours

Manual testing checklist:

- [ ] Test login flow
- [ ] Test warga list loading
- [ ] Test classification with auto-fill
- [ ] Test classification with manual entry
- [ ] Test training model workflow (both splits)
- [ ] Test model evaluation display
- [ ] Test role assignment modal
- [ ] Test welfare category filtering
- [ ] Test logout

**Acceptance:**
- All features work end-to-end
- No JavaScript errors
- Data persists correctly
- UI responsive

---

### Task 5.5: Desktop Application Build and Deployment Test
**Priority:** High | **Estimate:** 2 hours

- [ ] Build production executable: `wails build`
- [ ] Test executable on clean Windows machine
- [ ] Verify database auto-creation
- [ ] Test without Go installed
- [ ] Test application startup
- [ ] Test application shutdown
- [ ] Check for memory leaks

**Acceptance:**
- Executable runs on clean system
- No dependencies required (except WebView2)
- Database creates automatically
- No crashes or hangs

---

## Phase 6: Documentation and Delivery

### Task 6.1: Update User Manual
**Priority:** Medium | **Estimate:** 2 hours

Update `MANUAL_BOOK.md`:

- [ ] Add desktop application installation instructions
- [ ] Document new training model workflow
- [ ] Document model selection process
- [ ] Document role assignment feature
- [ ] Document filtering features
- [ ] Update screenshots

**Acceptance:**
- Manual covers all new features
- Screenshots current
- Instructions clear

---

### Task 6.2: Create Release Package
**Priority:** Medium | **Estimate:** 1 hour

- [ ] Build production executable
- [ ] Create installer (optional - Inno Setup)
- [ ] Include README with system requirements
- [ ] Include sample database (optional)
- [ ] Create release ZIP file

**Acceptance:**
- Release package ready for distribution
- All necessary files included
- README clear

---

## Task Summary

**Total Estimated Hours:** ~80 hours (~10 days)

**Critical Path:**
1. Wails setup and migration (Tasks 1.1-1.5)
2. Database migration (Tasks 2.1-2.2)
3. Backend services (Tasks 3.1-3.6)
4. Frontend enhancements (Tasks 4.1-4.5)
5. Excel verification testing (Task 5.1)
6. Integration testing (Tasks 5.2-5.5)

**Dependencies:**
- Task 1.x must complete before all others
- Task 2.x must complete before Task 3.x
- Task 3.x must complete before Task 4.x
- Task 4.x must complete before Task 5.x

**Risk Areas:**
- Excel calculation equivalence (highest priority)
- Wails learning curve if unfamiliar
- WebView2 compatibility on older Windows versions

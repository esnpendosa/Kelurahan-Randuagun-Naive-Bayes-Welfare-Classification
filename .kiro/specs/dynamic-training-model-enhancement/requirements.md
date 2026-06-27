# Requirements Document

## Introduction

This document specifies requirements for enhancing the existing Naive Bayes welfare classification system with dynamic training model capabilities, automatic data calculation, improved filtering, and streamlined user interface. The system currently classifies citizens into six welfare categories (Sangat Miskin, Miskin, Hampir Miskin, Rentan Miskin, Pas-pasan, Menengah ke Atas) using 36 indicators. This enhancement will enable users to work with multiple training models dynamically, automatically populate classification data, and filter results by specific welfare categories.

**CRITICAL PRIORITIES:**
1. **Calculation Accuracy**: All Naive Bayes calculations MUST match exactly with the Excel reference file ("klasifikasi naive bayes tambahan data.xlsx"). This is the highest priority requirement.
2. **Feature Completeness**: All requested features must be implemented exactly as specified in the user requirements.
3. **Desktop Application**: The application MUST run as a native Windows desktop application, not as a web browser-based application.

## Glossary

- **Classification_System**: The Naive Bayes-based desktop application that classifies citizen welfare levels
- **Desktop_Application**: A native Windows application using embedded webview (Wails/Lorca) or native UI framework, not requiring external web browser
- **Training_Model**: A trained Naive Bayes model using a specific dataset split
- **Training_Data**: Citizens marked as training samples (data_latih = 1 or data_latih_2 = 1)
- **Test_Data**: Citizens marked as test samples (data_latih = 0 or data_latih_2 = 0)
- **Indicator_Values**: The 36 welfare indicators (IM1-IM36) with categorical responses
- **Data_Split**: A partition strategy dividing citizens into training and test sets
- **Evaluation_Results**: Performance metrics including accuracy, precision, recall, F1-score, and confusion matrix
- **Automatic_Calculation**: System feature that populates indicator values from database when available
- **Dynamic_Model_Selection**: User ability to choose between multiple trained models for classification
- **Welfare_Category**: One of six classification levels (Sangat Miskin, Miskin, Hampir Miskin, Rentan Miskin, Pas-pasan, Menengah ke Atas)
- **Database**: SQLite database storing citizen data, indicators, and classification results
- **Webview**: Embedded browser component within desktop application for rendering UI without opening external browser

## Requirements

### Requirement 1: Automatic Data Calculation from Database

**User Story:** As an operator, I want indicator values to be automatically populated from the database when available, so that I don't have to manually re-enter data for existing citizens.

#### Acceptance Criteria

1. WHEN a citizen with existing indicator data is selected for classification, THE Classification_System SHALL automatically populate all 36 Indicator_Values from the Database
2. WHEN indicator data is not available for a selected citizen, THE Classification_System SHALL display empty form fields for manual entry
3. THE Classification_System SHALL retrieve Indicator_Values from the data_indikator table based on the warga_id
4. WHEN indicator values are auto-populated, THE Classification_System SHALL allow manual modification of any field before submission
5. THE Classification_System SHALL display a visual indication (e.g., pre-filled styling) to distinguish auto-populated fields from empty fields

### Requirement 2: Remove Training Data Labels from UI

**User Story:** As a user, I want training data labels removed from specific features, so that the interface is cleaner and less confusing.

#### Acceptance Criteria

1. THE Classification_System SHALL NOT display "data latih" (training data) labels in the Citizen_Data_Feature (warga.html page)
2. THE Classification_System SHALL NOT display training/test labels in the Classification_Feature dropdown selectors
3. WHEN displaying citizen lists, THE Classification_System SHALL show only essential information (NIK, name, address, welfare category)
4. THE Classification_System SHALL maintain internal tracking of training/test status without displaying it to end users
5. THE Training_Model_Feature SHALL continue to display training/test distinctions as this is essential for model evaluation

### Requirement 3: Welfare Category Filtering

**User Story:** As an administrator, I want to filter classification results by specific welfare categories, so that I can analyze "Sangat Miskin" and "Miskin" populations separately.

#### Acceptance Criteria

1. THE Classification_System SHALL provide a filter control in the Reports_Feature to select individual Welfare_Category values
2. WHEN a Welfare_Category filter is applied, THE Classification_System SHALL display only citizens matching the selected category
3. THE Classification_System SHALL support filtering by "Sangat Miskin" category individually
4. THE Classification_System SHALL support filtering by "Miskin" category individually
5. THE Classification_System SHALL support filtering by any of the six Welfare_Category values
6. WHEN no filter is applied, THE Classification_System SHALL display all classified citizens
7. THE Classification_System SHALL update distribution statistics to reflect filtered data only
8. THE Classification_System SHALL provide a clear filter control to reset to showing all categories

### Requirement 4: Dynamic Training Model Selection

**User Story:** As an administrator, I want to select between multiple training models before classification, so that I can compare model performance and choose the most accurate one.

#### Acceptance Criteria

1. THE Classification_System SHALL support at least two independent Data_Split configurations (Training 1/Test 1 and Training 2/Test 2)
2. WHEN viewing the Training_Model_Feature, THE Classification_System SHALL display data distribution for all available Data_Split configurations in separate tabs
3. WHEN the "Process Calculation" button is clicked, THE Classification_System SHALL prompt the user to select which Training_Model to use
4. THE Classification_System SHALL train the selected Training_Model using its corresponding Training_Data
5. THE Classification_System SHALL evaluate the selected Training_Model using its corresponding Test_Data
6. WHEN evaluation completes, THE Classification_System SHALL display Evaluation_Results including accuracy, precision, recall, F1-score, and confusion matrix
7. THE Classification_System SHALL clearly indicate which Training_Model was used for the displayed Evaluation_Results
8. THE Classification_System SHALL persist the selected Training_Model for subsequent classifications until changed by the user

### Requirement 5: Training Model Data Viewing

**User Story:** As an administrator, I want to view training and test data distributions without running evaluation, so that I can understand data splits before processing.

#### Acceptance Criteria

1. WHEN the Training_Model_Feature page loads, THE Classification_System SHALL display Training_Data count for Data_Split 1
2. WHEN the Training_Model_Feature page loads, THE Classification_System SHALL display Test_Data count for Data_Split 1
3. WHEN the Training_Model_Feature page loads, THE Classification_System SHALL display Training_Data count for Data_Split 2
4. WHEN the Training_Model_Feature page loads, THE Classification_System SHALL display Test_Data count for Data_Split 2
5. THE Classification_System SHALL display Welfare_Category distribution for each Data_Split as visual charts
6. THE Classification_System SHALL organize Data_Split displays in separate tabs for clarity
7. WHEN switching between tabs, THE Classification_System SHALL update displayed statistics without requiring page reload

### Requirement 6: Dynamic Data Role Assignment

**User Story:** As an administrator, I want to dynamically assign citizens to training or test roles for Data_Split 2, so that I can create custom dataset configurations.

#### Acceptance Criteria

1. THE Classification_System SHALL provide a dynamic assignment interface in the Training_Model_Feature
2. THE Classification_System SHALL display all citizens from the Database with their current Data_Split 2 role (training or test)
3. WHEN an administrator selects a citizen, THE Classification_System SHALL provide controls to change the citizen's role between training and test for Data_Split 2
4. WHEN a role change is saved, THE Classification_System SHALL update the data_latih_2 field in the warga table
5. THE Classification_System SHALL immediately reflect role changes in the Data_Split 2 distribution statistics
6. THE Classification_System SHALL preserve Data_Split 1 assignments independently from Data_Split 2 assignments
7. WHEN a citizen has no indicator data, THE Classification_System SHALL allow role assignment but indicate missing data status

### Requirement 7: New Classification Data Integration

**User Story:** As an operator, I want newly classified citizens to be automatically added to the training model feature dynamically, so that the system learns from new data.

#### Acceptance Criteria

1. WHEN a new citizen is classified, THE Classification_System SHALL store Indicator_Values in the data_indikator table
2. WHEN a new citizen is classified, THE Classification_System SHALL store the predicted Welfare_Category in the warga table label_kelas field
3. THE Classification_System SHALL make newly classified citizens available in the Dynamic_Data_Role_Assignment interface
4. WHEN viewing the Training_Model_Feature after new classifications, THE Classification_System SHALL include new citizens in the available data pool
5. THE Classification_System SHALL default new citizens to test role (data_latih_2 = 0) for Data_Split 2
6. THE Classification_System SHALL allow administrators to later reassign new citizens to training role if desired
7. WHEN new citizens are assigned to training role, THE Classification_System SHALL include them in subsequent model training for the selected Data_Split

### Requirement 8: Pre-filled Classification for Existing Citizens

**User Story:** As an operator, I want to classify existing citizens with a single click, so that I can quickly process citizens who already have indicator data.

#### Acceptance Criteria

1. WHEN selecting a citizen with complete indicator data from the classification form, THE Classification_System SHALL auto-populate all 36 Indicator_Values
2. THE Classification_System SHALL enable the "Process Classification" button when all required fields are populated
3. WHEN the "Process Classification" button is clicked with pre-filled data, THE Classification_System SHALL use the currently active Training_Model for prediction
4. THE Classification_System SHALL update the citizen's Welfare_Category in the Database
5. WHEN classification completes, THE Classification_System SHALL redirect to the results page showing the predicted category and probability distribution
6. THE Classification_System SHALL store the full probability distribution as JSON in the hasil_klasifikasi table

### Requirement 9: Streamlined Authentication Flow

**User Story:** As a user, I want a single logout option instead of redundant logout and exit buttons, so that the interface is simpler and less confusing.

#### Acceptance Criteria

1. THE Classification_System SHALL provide a single "Logout" button in the navigation interface
2. WHEN the "Logout" button is clicked, THE Classification_System SHALL clear the user session and redirect to the login page
3. THE Classification_System SHALL remove duplicate logout and exit application controls from the user interface
4. THE Classification_System SHALL maintain a separate "Exit Application" button only in the desktop executable version (if applicable)
5. THE Classification_System SHALL preserve session security by invalidating session tokens on logout

### Requirement 10: Training Model Performance Evaluation

**User Story:** As an administrator, I want to see detailed performance metrics for each training model, so that I can determine which model performs better for my dataset.

**CRITICAL**: All calculation results (accuracy, precision, recall, F1-score, confusion matrix) MUST match exactly with the calculations in the Excel file "klasifikasi naive bayes tambahan data.xlsx". Any deviation is considered a critical bug.

#### Acceptance Criteria

1. WHEN model evaluation completes, THE Classification_System SHALL calculate accuracy as (correct predictions / total predictions) matching Excel calculations
2. WHEN model evaluation completes, THE Classification_System SHALL calculate precision for each Welfare_Category matching Excel formulas
3. WHEN model evaluation completes, THE Classification_System SHALL calculate recall for each Welfare_Category matching Excel formulas
4. WHEN model evaluation completes, THE Classification_System SHALL calculate F1-score for each Welfare_Category matching Excel formulas
5. THE Classification_System SHALL generate a 6x6 confusion matrix showing predicted vs actual Welfare_Category for all Test_Data matching Excel results
6. THE Classification_System SHALL display macro-averaged precision, recall, and F1-score across all categories matching Excel calculations
7. WHEN a Welfare_Category filter is applied, THE Classification_System SHALL recalculate metrics for only the filtered category
8. THE Classification_System SHALL display a timestamp indicating when the model was last trained
9. THE Classification_System SHALL use Laplace Smoothing with the same parameters as Excel calculations
10. WHEN comparing system results with Excel, THE Classification_System SHALL produce identical probability distributions (within 0.001% tolerance)

### Requirement 11: Database Schema Extensions

**User Story:** As a system, I need extended database schema to support multiple data splits, so that dynamic training model features can function correctly.

#### Acceptance Criteria

1. THE Classification_System SHALL add a data_latih_2 column to the warga table with integer type and default value 0
2. THE data_latih_2 column SHALL indicate training status (1) or test status (0) for Data_Split 2
3. THE Classification_System SHALL preserve existing data_latih column functionality for Data_Split 1
4. THE Classification_System SHALL support querying citizens by data_latih_2 values for Data_Split 2 operations
5. WHEN a database migration is needed, THE Classification_System SHALL execute ALTER TABLE commands safely
6. THE Classification_System SHALL maintain backward compatibility with existing data where data_latih_2 may be null or missing

### Requirement 12: Filter Persistence and Reset

**User Story:** As a user, I want filter selections to persist during my session and a clear way to reset them, so that I can work efficiently with filtered data views.

#### Acceptance Criteria

1. WHEN a Welfare_Category filter is applied in the Reports_Feature, THE Classification_System SHALL maintain the filter selection when navigating within the same session
2. WHEN a Welfare_Category filter is applied in the Training_Model_Feature evaluation, THE Classification_System SHALL maintain the filter for the duration of viewing results
3. THE Classification_System SHALL provide a "Clear Filter" or "Show All" button to remove active filters
4. WHEN the "Clear Filter" button is clicked, THE Classification_System SHALL display all data and reset filter controls to default state
5. THE Classification_System SHALL visually indicate when a filter is active (e.g., highlighted filter button, active filter badge)
6. WHEN switching between Training_Model Data_Split tabs, THE Classification_System SHALL reset filters to avoid confusion



### Requirement 13: Excel Calculation Verification

**User Story:** As a quality assurance tester, I want to verify that system calculations match Excel reference data exactly, so that I can ensure calculation accuracy.

**CRITICAL REQUIREMENT**: This is the highest priority requirement. All implementation must be validated against Excel.

#### Acceptance Criteria

1. THE Classification_System SHALL load training and test data from the same data splits as Excel sheets:
   - "Data Training 1" sheet → data_latih = 1
   - "Data Uji 1" sheet → data_latih = 0
   - "Data Training 2" sheet → data_latih_2 = 1
   - "Data Uji 2" sheet → data_latih_2 = 0
2. WHEN calculating prior probabilities P(C), THE Classification_System SHALL produce identical values to Excel calculations
3. WHEN calculating likelihood P(X|C) with Laplace Smoothing, THE Classification_System SHALL produce identical values to Excel calculations
4. WHEN classifying a test citizen, THE Classification_System SHALL produce the same predicted class as Excel
5. WHEN classifying a test citizen, THE Classification_System SHALL produce probability distributions within 0.001% of Excel values
6. WHEN calculating confusion matrix, THE Classification_System SHALL produce identical counts to Excel evaluation results
7. WHEN calculating accuracy, precision, recall, and F1-score, THE Classification_System SHALL produce values within 0.01% of Excel metrics
8. THE Classification_System SHALL provide a verification mode that compares system results against Excel for validation
9. WHEN verification detects discrepancies, THE Classification_System SHALL log detailed comparison reports showing differences
10. THE Classification_System SHALL use the same Laplace Smoothing formula: P(Xi=v|Ck) = (count + 1) / (total_in_class + |V|) where |V| is the number of unique values for feature Xi


### Requirement 14: Native Windows Desktop Application

**User Story:** As a user, I want to run the application as a native Windows desktop application without opening a web browser, so that it feels like a proper desktop program.

**CRITICAL REQUIREMENT**: The application must run as a standalone Windows executable without requiring external browser.

#### Acceptance Criteria

1. THE Classification_System SHALL be packaged as a standalone Windows executable (.exe file)
2. WHEN the user double-clicks the executable, THE Classification_System SHALL launch in its own application window, NOT in an external web browser
3. THE Classification_System SHALL use an embedded webview component (such as Wails or Lorca) to render the HTML/CSS/JavaScript UI
4. THE Classification_System SHALL display a custom application icon in the title bar and taskbar
5. THE Classification_System SHALL have a custom application title (e.g., "Sistem Klasifikasi Kesejahteraan Randuagung")
6. WHEN the application starts, THE Classification_System SHALL display its own window with configurable size (minimum 1024x768 pixels)
7. THE Classification_System SHALL NOT require users to manually open a browser and navigate to localhost URL
8. THE Classification_System SHALL automatically bind the Go backend to the webview frontend
9. WHEN the user closes the application window, THE Classification_System SHALL cleanly shutdown the backend server and database connections
10. THE Classification_System SHALL support Windows 10 and Windows 11 operating systems
11. THE Classification_System SHALL bundle all necessary dependencies (no external browser installation required)
12. THE Classification_System SHALL maintain the same HTML/CSS/JavaScript frontend code structure for UI rendering
13. WHEN the application starts, THE Classification_System SHALL check if the database file exists and create it if missing
14. THE Classification_System SHALL support standard desktop window controls (minimize, maximize, close)
15. THE Classification_System SHALL optionally support fullscreen mode for kiosk-style deployment

#### Technical Notes

**Recommended Technology Stack:**
- **Wails v2** (preferred): Modern Go + Web frontend framework with native Windows webview
  - Uses Microsoft Edge WebView2 on Windows
  - Single executable output
  - Hot reload during development
  - Native Go ↔ JavaScript bindings
  
- **Lorca** (alternative): Simpler Go + Chrome/Chromium integration
  - Uses installed Chrome/Edge browser engine
  - Lighter weight than Wails
  - Less features but easier setup

**Migration Path:**
1. Keep existing Go backend code (main.go, classifier, db packages)
2. Keep existing HTML/CSS/JavaScript frontend files
3. Integrate with Wails/Lorca to create desktop wrapper
4. Replace http.ListenAndServe with webview window creation
5. Add native bindings for Go functions called from JavaScript
6. Build as Windows executable with embedded frontend assets

**Build Output:**
- Single .exe file (e.g., `Klasifikasi-Warga-Randuagung.exe`)
- Embedded frontend assets (HTML, CSS, JS, images)
- Database file (data_skripsi.db) in same directory or AppData folder

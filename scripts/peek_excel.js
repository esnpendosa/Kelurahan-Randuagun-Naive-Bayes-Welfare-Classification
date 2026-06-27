import XLSX from 'xlsx';

const workbook = XLSX.readFile('klasifikasi naive bayes tambahan data.xlsx');

function peekSheetNonEmpty(sheetName) {
  console.log(`\n--- PEEK SHEET: ${sheetName} ---`);
  const sheet = workbook.Sheets[sheetName];
  if (!sheet) {
    console.log('Sheet not found!');
    return;
  }
  const rows = XLSX.utils.sheet_to_json(sheet, { header: 1 });
  const nonEmptyRows = rows.filter(r => r && r[1] !== undefined && String(r[1]).trim() !== '' && r[1] !== 'Kepala Rumah Tangga');
  console.log('Total parsed rows:', rows.length);
  console.log('Non-empty rows count:', nonEmptyRows.length);
  if (nonEmptyRows.length > 0) {
    console.log('First non-empty row:', nonEmptyRows[0]);
    console.log('Last non-empty row:', nonEmptyRows[nonEmptyRows.length - 1]);
  }
}

peekSheetNonEmpty('Data Training 1');
peekSheetNonEmpty('Data Uji 1');
peekSheetNonEmpty('Data Training 2');
peekSheetNonEmpty('Data Uji 2');

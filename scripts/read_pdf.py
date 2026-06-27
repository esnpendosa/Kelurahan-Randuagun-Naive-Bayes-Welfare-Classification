import sys
import os

try:
    import pypdf
except ImportError:
    import subprocess
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pypdf"])
    import pypdf

def main():
    pdf_path = "proposal skripsi fixxx TABARAKALLAH.pdf"
    if not os.path.exists(pdf_path):
        print("PDF not found")
        return
    
    reader = pypdf.PdfReader(pdf_path)
    print(f"Total pages: {len(reader.pages)}")
    
    # Search for keywords
    keywords = ["akurasi", "confusion", "matriks", "uji", "training", "validasi"]
    for i, page in enumerate(reader.pages):
        text = page.extract_text()
        found = [kw for kw in keywords if kw.lower() in text.lower()]
        if found:
            print(f"Page {i+1} has keywords: {found}")
            # print snippet of text if matches confusion matrix
            if "confusion" in text.lower() or "akurasi" in text.lower():
                print(f"--- Page {i+1} Text Snippet ---")
                lines = text.split("\n")
                for line in lines:
                    if any(kw in line.lower() for kw in ["confusion", "akurasi", "miskin", "uji", "latih"]):
                        print(line[:100])

if __name__ == "__main__":
    main()

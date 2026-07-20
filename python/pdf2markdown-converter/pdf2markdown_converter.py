import pymupdf4llm
from pathlib import Path


def convert_pdf_pymupdf(
    pdf_path: str | Path, output_dir: str | Path = "./output_pymupdf"
) -> Path:
    """Converts a PDF to Markdown and extracts image using PyMuPDF4LLM.

    Args:
        pdf_path: Path to the input PDF file.
        output_dir: Folder to store the resulting .md file and images.

    Returns:
        Path to the written Markdown file.
    """
    pdf_path = Path(pdf_path)
    output_dir = Path(output_dir)
    images_dir = output_dir / "images"

    output_dir.mkdir(parents=True, exist_ok=True)
    images_dir.mkdir(parents=True, exist_ok=True)

    markdown_content = pymupdf4llm.to_markdown(
        doc=pdf_path,
        write_images=True,
        image_path=str(images_dir),
        images_format="png",
        dpi=150,
        header=False,
        footer=False,
    )

    output_md_path = output_dir / f"{pdf_path.stem}.md"
    output_md_path.write_text(markdown_content, "utf-8")

    print(f"Done! Markdown saved to: {output_md_path}")

    return output_md_path


def main():
    pdf_file = "sample.pdf"
    convert_pdf_pymupdf(pdf_path=pdf_file, output_dir="./output_result")

if __name__ == "__main__":
    main()

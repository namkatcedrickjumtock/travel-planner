# how to run the script:
# python3 compress_image.py image.png
# python3 compress_image.py image.png --remove-bg
# python3 compress_image.py image.png --remove-bg --max-width 1000
 
import argparse
import os
import io
from PIL import Image
from rembg import remove

# Disable Pillow pixel limit to avoid decompression bomb errors
Image.MAX_IMAGE_PIXELS = None

def compress_image(input_path, output_path=None, quality=70, max_width=1920, max_height=1920, remove_bg=False):
    try:
        # Check input file exists
        if not os.path.exists(input_path):
            print(f"Error: Input file not found: {input_path}")
            return

        # Load image
        with open(input_path, "rb") as f:
            input_data = f.read()

        # Remove background if requested
        if remove_bg:
            output_data = remove(input_data)
            img = Image.open(io.BytesIO(output_data)).convert("RGBA")
        else:
            img = Image.open(input_path)

        # Resize while keeping aspect ratio
        img.thumbnail((max_width, max_height))

        # Default output path if not provided
        if not output_path:
            base, ext = os.path.splitext(input_path)
            if remove_bg:
                output_path = f"{base}_no_bg.png"
            else:
                output_path = f"{base}_compressed.jpg"

        # Save image
        if remove_bg:
            # Preserve transparency
            img.save(output_path, "PNG", optimize=True)
        else:
            # Convert to RGB if necessary
            if img.mode in ("RGBA", "P"):
                img = img.convert("RGB")
            img.save(output_path, "JPEG", quality=quality, optimize=True)

        # Logging
        original_size = os.path.getsize(input_path) / 1024
        compressed_size = os.path.getsize(output_path) / 1024
        print(f"Output saved to: {output_path}")
        print(f"Original size: {original_size:.2f} KB")
        print(f"Compressed size: {compressed_size:.2f} KB")

    except Exception as e:
        print(f"Error processing {input_path}: {e}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Compress, resize, and optionally remove background from images")

    parser.add_argument("input", help="Path to input image")
    parser.add_argument("-o", "--output", help="Output path (optional)")
    parser.add_argument("-q", "--quality", type=int, default=70, help="JPEG quality (1-95)")
    parser.add_argument("--max-width", type=int, default=1920, help="Maximum width")
    parser.add_argument("--max-height", type=int, default=1920, help="Maximum height")
    parser.add_argument("--remove-bg", action="store_true", help="Remove background (produces PNG)")

    args = parser.parse_args()

    compress_image(
        input_path=args.input,
        output_path=args.output,
        quality=args.quality,
        max_width=args.max_width,
        max_height=args.max_height,
        remove_bg=args.remove_bg
    ) 
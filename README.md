# merge2pdf - Simplest tool for merging Images and PDFs into PDF

Merge Image and PDF files (optionally with selective pages) with lossless quality.
It will not convert PDF pages (with texts, images, forms) into flat image, everything will remain as is.

### Install

It's just a single binary file, no external dependencies. 
Just download the appropriate version of [executable from latest release](https://github.com/ajaxray/merge2pdf/releases) for your OS. Then rename and give it execute permission.
```bash
mv merge2pdf_linux-amd64 merge2pdf  
sudo chmod +x merge2pdf
```

If you want to install it globally (run from any directory of your system), put it in your systems $PATH directory.
```bash
sudo mv merge2pdf /usr/local/bin/merge2pdf
```
Done! 

### How to use

```bash
# Merge multiple PDFs
merge2pdf output.pdf input1.pdf input2.pdf path/to/other.pdf ...

# Merge 1st page of input1.pdf, full input2.pdf and 2nd, 3rd, 4th page of input3.pdf  
merge2pdf output.pdf input1.pdf~1 input2.pdf input3.pdf~2,3,4

# Merge multiple Images
merge2pdf output.pdf image1.jpg image2.jpg path/to/other.png ...

# Mixing up PDF, PDF Pages and Images
merge2pdf output.pdf doc1.pdf~1,2 image1.jpg image2.png path/to/other.pdf ...
```

### Fine tuning Image pages

When adding Images, by default the size for image containing pages will be same to image size with 1 inch margin on each side. But you may set custom margins and resize to standard Print size.
```bash
# Set image page size to A4. Other Options are A3, Legal and Letter
merge2pdf -s A4 output.pdf image1.jpg image2.jpg

# Set image page size and margin (left, right, top, bottom).
merge2pdf -s A3 -m .5,.5,1,1 output.pdf image1.jpg image2.jpg
# margin can be set without size

# Scale images to page width, with aspect ratio
# To scale to height, use --scale-height
merge2pdf -s A3 -m .5,.5,1,1 --scale-width output.pdf image1.jpg image2.jpg
```
_Note: When adding images and PDFs together, these options will effect ONLY Image pages._


If your filename contains space or [some special characters](https://unix.stackexchange.com/a/270979), 
then quote the filepaths along with page numbers. For safety, you can quote them always. 
```bash
merge2pdf output.pdf "With Space.pdf" "without-space.pdf" "with space and pages.pdf~2,3,4"
```

To see details of options, 
```bash
merge2pdf -h
```

### Roadmap

✅ Merge multiple PDFs without loosing quality  
✅ Merge multiple PDFs with **selective pages**  
✅ Adding Images  
✅ Mixing up Images and PDFs  
◻️ Merge all (Image and PDF) files from directory  
✅ Option to Resize Images to reduce filesize  
◻️ Option to Greyscale Images to reduce filesize  
◻️ Option to set files and pages as JSON config to make usages from other app more convenient  

### Contribute

If you fix a bug or want to add/improve a feature, 
and it's alligned with the focus (merging with ease) of this tool, 
I will be glad to accept your PR. :) 

### Thanks

This tool was made using the beautiful [Unidoc](https://unidoc.io/) library. Thanks and ❤️ to **Unidoc**.

---
> "This is the Book about which there is no doubt, a guidance for those conscious of Allah" - [Al-Quran](http://quran.com)

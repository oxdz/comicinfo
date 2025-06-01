package shonenmagazine

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

func _drawImage(
	dc *gg.Context, // 目标画布
	srcImage image.Image, // 源图像
	srcX, srcY int, // 源图像裁剪区域左上角坐标
	srcWidth, srcHeight int, // 源图像裁剪区域尺寸
	dstX, dstY int, // 目标位置
	dstWidth, dstHeight int, // 目标尺寸
) {
	// 创建裁剪区域的子图像
	srcBounds := srcImage.Bounds()
	subImage := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))

	// 复制源图像的指定区域到子图像
	for y := 0; y < srcHeight && srcY+y < srcBounds.Max.Y; y++ {
		for x := 0; x < srcWidth && srcX+x < srcBounds.Max.X; x++ {
			subImage.Set(x, y, srcImage.At(srcX+x, srcY+y))
		}
	}

	// 绘制子图像到目标位置
	dc.DrawImage(subImage, dstX, dstY)

	// 如果需要缩放，可以使用以下方式
	if srcWidth != dstWidth || srcHeight != dstHeight {
		// 创建缩放后的图像
		scaledImage := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
		draw.BiLinear.Scale(scaledImage, scaledImage.Bounds(), subImage, subImage.Bounds(), draw.Over, nil)
		dc.DrawImage(scaledImage, dstX, dstY)
	}
}

func DrawImage(srcImage image.Image, dst string, seed int) error {
	mappings := le(seed)

	bd := srcImage.Bounds()
	x := bd.Dx()
	y := bd.Dy()

	// 创建目标画布
	dc := gg.NewContextForImage(srcImage)

	// 假设我们有一个映射，表示如何重新排列图像块

	// 块大小
	v := BlockSize(x, y, 4)
	blockWidth, blockHeight := v.Width, v.Height
	// 应用映射
	for _, m := range mappings {
		_drawImage(
			dc,
			srcImage,
			m.Source.X*blockWidth, m.Source.Y*blockHeight,
			blockWidth, blockHeight,
			m.Dest.X*blockWidth, m.Dest.Y*blockHeight,
			blockWidth, blockHeight,
		)
	}

	// 保存结果
	if err := dc.SavePNG(dst); err != nil {
		fmt.Println("Error saving image:", err)
	}

	return nil
}

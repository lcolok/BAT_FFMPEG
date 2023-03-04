package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/atotto/clipboard"
)

func main() {
	var inputDir string

	// 检查是否有命令行参数传递进来
	if len(os.Args) > 1 {
		// 如果有，尝试将其作为输入目录路径
		inputDir = os.Args[1]
	} else {
		// 否则，从剪贴板中获取输入目录
		inputDirFromClipboard, err := clipboard.ReadAll()
		if err != nil {
			panic(err)
		}
		inputDir = inputDirFromClipboard
	}

	// 如果没有输入目录，尝试从拖放的文件夹中获取
	if inputDir == "" && len(os.Args) == 1 {
		if len(os.Args) == 1 {
			// 获取文件所在目录
			execPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

			if filepath.Ext(execPath) == ".exe" {
				// 如果是.exe文件，则尝试获取拖入的文件夹路径
				if len(os.Args) > 1 {
					inputDir = os.Args[1]
				} else {
					// 如果没有命令行参数，提示用户将文件夹拖入.exe文件上
					fmt.Println("请将文件夹拖入此可执行文件上以继续")
					return
				}
			} else {
				// 如果是.go文件，则尝试获取拖入的文件夹路径
				if len(os.Args) > 2 {
					inputDir = os.Args[2]
				} else {
					// 如果没有命令行参数，提示用户将文件夹拖入.exe文件上
					fmt.Println("请将文件夹拖入此可执行文件上以继续")
					return
				}
			}
			inputDir = filepath.Clean(inputDir)
		}
	}

	// 检查输入目录是否是一个合法的目录路径
	fileInfo, err := os.Stat(inputDir)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		panic("输入目录不存在或不是一个目录路径")
	}

	// 获取输入目录的父级目录路径
	parentDir := filepath.Dir(inputDir)

	// 获取输入目录中的所有PNG文件名
	pattern := filepath.Join(inputDir, "*.png")
	pngFiles, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	// 解析文件名以确定文件名规律
	re := regexp.MustCompile(`(\d+)\.png$`)
	var startFrame int
	var frameRate int
	for _, filename := range pngFiles {
		matches := re.FindStringSubmatch(filename)
		if matches != nil {
			frameNumber, err := strconv.Atoi(matches[1])
			if err != nil {
				panic(err)
			}
			if startFrame == 0 || frameNumber < startFrame {
				startFrame = frameNumber
			}
			frameRate++
		}
	}
	if frameRate == 0 {
		panic("未找到PNG文件")
	}

	// 解析文件名以确定数字的位数
	re2 := regexp.MustCompile(`(\d+)\.png$`)
	var numDigits int
	for _, filename := range pngFiles {
		matches := re2.FindStringSubmatch(filename)
		if matches != nil {
			num := len(matches[1])
			if numDigits == 0 || num > numDigits {
				numDigits = num
			}
		}
	}
	if numDigits == 0 {
		panic("未找到PNG文件")
	}

	// 使用解析出的文件名规律构建FFmpeg命令行参数
	args := []string{
		"-start_number", strconv.Itoa(startFrame),
		"-framerate", strconv.Itoa(frameRate),
		"-i", inputDir + "/" + re.ReplaceAllString(filepath.Base(pngFiles[0]), fmt.Sprintf("%%%dd.png", numDigits)),
		"-c:v", "libx264", "-y",
		filepath.Join(parentDir, "output.mp4"),
	}

	// 执行FFmpeg命令
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("输出文件已保存为 %s\n", filepath.Join(parentDir, "output.mp4"))
}

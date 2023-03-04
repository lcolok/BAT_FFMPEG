package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	inputDir := getInputDir()

	validateInputDir(inputDir)

	parentDir := getParentDir(inputDir)

	pngFiles := getPNGFiles(inputDir)

	startFrame, _ := parsePNGFiles(pngFiles)

	defaultFrameRate := 30

	frameRate := getUserFrameRate(defaultFrameRate)

	numDigits := getNumDigits(pngFiles)

	outputPath := filepath.Join(parentDir, "output.mp4")
	buildFFmpegCommand(inputDir, startFrame, frameRate, numDigits, outputPath, pngFiles)

	fmt.Printf("输出文件已保存为 %s\n", outputPath)
}

func getUserFrameRate(defaultFrameRate int) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入帧速率（默认为30）：")
	frameRateString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取帧速率错误：", err)
		os.Exit(1)
	}
	frameRateString = strings.TrimSpace(frameRateString)

	var frameRate int
	if frameRateString == "" {
		frameRate = defaultFrameRate
	} else {
		frameRate, err = strconv.Atoi(frameRateString)
		if err != nil {
			fmt.Println("帧速率格式错误：", err)
			os.Exit(1)
		}
	}
	return frameRate
}

func getInputDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	inputDirFromClipboard, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	return inputDirFromClipboard
}

func validateInputDir(inputDir string) {
	fileInfo, err := os.Stat(inputDir)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		panic("输入目录不存在或不是一个目录路径")
	}
}

func getParentDir(inputDir string) string {
	return filepath.Dir(inputDir)
}

func getPNGFiles(inputDir string) []string {
	pattern := filepath.Join(inputDir, "*.png")
	pngFiles, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	if len(pngFiles) == 0 {
		panic("未找到PNG文件")
	}
	return pngFiles
}

func parsePNGFiles(pngFiles []string) (int, int) {
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
	return startFrame, frameRate
}

func getNumDigits(pngFiles []string) int {
	re := regexp.MustCompile(`(\d+)\.png$`)
	var numDigits int
	for _, filename := range pngFiles {
		matches := re.FindStringSubmatch(filename)
		if matches != nil {
			num := len(matches[1])
			if numDigits == 0 || num > numDigits {
				numDigits = num
			}
		}
	}
	return numDigits
}

func buildFFmpegCommand(inputDir string, startFrame, frameRate, numDigits int, outputPath string, pngFiles []string) {
	// 通过正则表达式匹配文件名以确定文件名规律
	fileNameRegex := regexp.MustCompile(`(\d+)\.png$`)
	// 生成FFmpeg命令的参数
	// 获取PNG文件的基础名称（不带目录路径和文件扩展名）
	baseName := filepath.Base(pngFiles[0])
	// 用正则表达式替换基础文件名中的数字以生成文件名格式字符串
	baseName = fileNameRegex.ReplaceAllString(baseName, fmt.Sprintf("%%%dd.png", numDigits))
	// 组合输入目录和格式化的文件名以生成输入路径
	inputPath := filepath.Join(inputDir, baseName)
	// 生成FFmpeg命令的参数
	ffmpegArgs := []string{
		"-start_number", strconv.Itoa(startFrame), // 设置开始帧数
		"-framerate", strconv.Itoa(frameRate), // 设置帧率
		"-i", inputPath, // 设置输入路径
		"-c:v", "libx264", // 设置视频编解码器为libx264
		"-y",       // 允许覆盖输出文件
		outputPath, // 设置输出路径
	}
	// 创建一个FFmpeg命令对象
	ffmpegCmd := exec.Command("ffmpeg", ffmpegArgs...)
	// 将输出和错误输出流重定向到标准输出和标准错误输出流
	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr
	// 设置FFmpeg命令执行的工作目录
	ffmpegCmd.Dir = inputDir
	// 执行FFmpeg命令
	if err := ffmpegCmd.Run(); err != nil {
		// 如果出现错误，则中止程序
		panic(err)
	}
}

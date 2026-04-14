package main

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
)

type SinhVien struct {
	Masv string `json:"masv"`
	Ten  string `json:"ten"`
	Tuoi int    `json:"tuoi"`
}

type MonHoc struct {
	Mamon  string `json:"mamon"`
	Ten    string `json:"ten"`
	Tinchi int    `json:"tinchi"`
}

type DangKy struct {
	Masv  string  `json:"masv"`
	Mamon string  `json:"mamon"`
	Diem  float64 `json:"diem"`
}

type DiemDanh struct {
	Masv  string `json:"masv"`
	Mamon string `json:"mamon"`
	Buoi  int    `json:"buoi"`
	CoMat bool   `json:"comat"`
}

type LMS struct {
	SinhViens []SinhVien `json:"sinhviens"`
	MonHocs   []MonHoc   `json:"monhocs"`
	DangKys   []DangKy   `json:"dangkys"`
	DiemDanhs []DiemDanh `json:"diemdanhs"`
}

var data LMS

func save() {
	file, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("data.json", file, 0644)
}

func load() {
	file, err := os.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(file, &data)
	}
}

func tinhGPA(masv string) float64 {
	var tong float64
	var tin int

	for _, dk := range data.DangKys {
		if dk.Masv == masv {
			for _, m := range data.MonHocs {
				if m.Mamon == dk.Mamon {
					tong += dk.Diem * float64(m.Tinchi)
					tin += m.Tinchi
				}
			}
		}
	}

	if tin == 0 {
		return 0
	}
	return tong / float64(tin)
}

func main() {
	load()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "LMS API is running",
			"endpoints": gin.H{
				"GET /sinhvien":        "Xem danh sach sinh vien",
				"POST /sinhvien":       "Them sinh vien",
				"DELETE /sinhvien/:id": "Xoa sinh vien",
				"GET /mon":             "Xem mon hoc",
				"POST /mon":            "Them mon hoc",
				"POST /dangky":         "Dang ky mon hoc",
				"DELETE /dangky":       "Huy dang ky mon hoc",
				"PUT /dangky":          "Xem mon da dang ky",
				"GET /gpa/:masv":       "Xem GPA",
				"POST /diemdanh":       "Điem danh",
				"GET /diemdanh/:mamon": "Xem danh sach diem danh theo mon",
				"DELETE /diemdanh":     "Xoa diem danh",
			},
			"how_to_test": "Su dung Postman hoac curl de kiem tra API",
		})
	})

	r.POST("/sinhvien", func(c *gin.Context) {
		var sv SinhVien
		c.BindJSON(&sv)

		data.SinhViens = append(data.SinhViens, sv)
		save()

		c.JSON(200, sv)
	})

	r.GET("/sinhvien", func(c *gin.Context) {
		c.JSON(200, data.SinhViens)
	})

	r.DELETE("/sinhvien/:id", func(c *gin.Context) {
		id := c.Param("id")

		for i, sv := range data.SinhViens {
			if sv.Masv == id {
				data.SinhViens = append(data.SinhViens[:i], data.SinhViens[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.POST("/mon", func(c *gin.Context) {
		var m MonHoc
		c.BindJSON(&m)

		data.MonHocs = append(data.MonHocs, m)
		save()

		c.JSON(200, m)
	})

	r.GET("/mon", func(c *gin.Context) {
		c.JSON(200, data.MonHocs)
	})

	r.POST("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		data.DangKys = append(data.DangKys, dk)
		save()

		c.JSON(200, dk)
	})

	r.PUT("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		for i := range data.DangKys {
			if data.DangKys[i].Masv == dk.Masv && data.DangKys[i].Mamon == dk.Mamon {
				data.DangKys[i].Diem = dk.Diem
				save()
				c.JSON(200, gin.H{"msg": "updated"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.DELETE("/dangky", func(c *gin.Context) {
		var dk DangKy
		c.BindJSON(&dk)

		for i := range data.DangKys {
			if data.DangKys[i].Masv == dk.Masv && data.DangKys[i].Mamon == dk.Mamon {
				data.DangKys = append(data.DangKys[:i], data.DangKys[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.GET("/gpa/:masv", func(c *gin.Context) {
		id := c.Param("masv")
		c.JSON(200, gin.H{"gpa": tinhGPA(id)})
	})

	r.POST("/diemdanh", func(c *gin.Context) {
		var dd DiemDanh
		c.BindJSON(&dd)

		data.DiemDanhs = append(data.DiemDanhs, dd)
		save()

		c.JSON(200, dd)
	})

	r.GET("/diemdanh/:mamon", func(c *gin.Context) {
		mamon := c.Param("mamon")

		var result []DiemDanh
		for _, d := range data.DiemDanhs {
			if d.Mamon == mamon {
				result = append(result, d)
			}
		}

		c.JSON(200, result)
	})

	r.DELETE("/diemdanh", func(c *gin.Context) {
		var dd DiemDanh
		c.BindJSON(&dd)

		for i, d := range data.DiemDanhs {
			if d.Masv == dd.Masv && d.Mamon == dd.Mamon && d.Buoi == dd.Buoi {
				data.DiemDanhs = append(data.DiemDanhs[:i], data.DiemDanhs[i+1:]...)
				save()
				c.JSON(200, gin.H{"msg": "deleted"})
				return
			}
		}

		c.JSON(404, gin.H{"error": "not found"})
	})

	r.Run(":8080")
}

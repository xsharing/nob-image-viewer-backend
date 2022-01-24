package design

import . "goa.design/goa/v3/dsl"

var _ = API("image-viewer", func() {
	Title("Image Viewer Service")
	Description("HTTP service for viewing image files hosted on S3")
	Server("image-viewer", func() {
		Host("localhost", func() { URI("http://localhost:8088") })
	})
})

var Image = Type("Image", func() {
	Attribute("name", String, "name")
	Attribute("thumbnail_url", String, "signed thumbnail image url")
	Attribute("original_url", String, "signed original image url")
})

var _ = Service("image-viewer", func() {
	Description("The image-viewer service")
	// Method describes a service method (endpoint)
	Method("folders", func() {
		Result(ArrayOf(String))
		HTTP(func() {
			GET("/folders")
			Response(StatusOK, func() {
				ContentType("application/json")
			})
		})
	})

	Method("images", func() {
		Payload(func() {
			Attribute("folder", String)
		})
		Result(ArrayOf(Image))
		HTTP(func() {
			GET("/images")
			Param("folder")
			Response(StatusOK, func() {
				ContentType("application/json")
			})
		})
	})
})

package adaptertemplate

import (
	"fmt"

	"go.uber.org/fx"
)

// AsRoute là helper function để annotate controller constructor vào Fx group
//
// Parameters:
//   - f: Constructor function (ví dụ: NewUserController)
//   - groupTag: Tên Fx group (ví dụ: "dynamicControllers", "apiRoutes")
//   - annotation: Các fx.Annotation tùy chọn khác
//
// Returns: Annotated function để dùng với fx.Provide
//
// Example:
//
//	fx.Provide(
//	    AsRoute(NewProductController, "dynamicControllers"),
//	)
//
//	// Constructor signature: func NewProductController(router *gin.Engine) *ProductController
//	// Constructor phải return type implement ICoreController interface
func AsRoute(f any, groupTag string, annotation ...fx.Annotation) any {
	annotation = append(annotation,
		fx.As(new(ICoreController)),
		fx.ResultTags(fmt.Sprintf(`group:"%s"`, groupTag)))
	return fx.Annotate(f, annotation...)
}

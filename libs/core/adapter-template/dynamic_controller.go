package adaptertemplate

// ICoreController là marker interface cho controllers sử dụng dynamic registration
// Controllers implement interface này sẽ có tất cả methods với signature func(context.Context)
// được tự động gọi thông qua reflection
type ICoreController interface {
	// Marker interface - không có methods
}

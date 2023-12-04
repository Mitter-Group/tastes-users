package instagram

// InstagramProperties implementa ProviderProperties para el proveedor Instagram
type InstagramCallbackRequestProperties struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
}

// Implementa los métodos de la interfaz ProviderProperties
func (i InstagramCallbackRequestProperties) GetCode() string {
	return i.Code
}

func (i InstagramCallbackRequestProperties) GetState() string {
	// Instagram no proporciona 'state', puedes retornar un valor predeterminado o manejarlo según tus necesidades
	return "N/A"
}

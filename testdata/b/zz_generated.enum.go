package b

func (PullPolicy) OpenAPISchemaEnum() []interface{} {
	return []interface{}{
		PullAlways,
		PullIfNotPresent,
		PullNever,
	}
}

package decoder

func (d *Decoder) ScanNull() error {
	c := d.char()
	if c != 'n' {
		return d.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	d.pos++
	c = d.char()
	if c != 'u' {
		return d.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	d.pos++
	c = d.char()
	if c != 'l' {
		return d.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	d.pos++
	c = d.char()
	if c != 'l' {
		return d.ErrorF("[Blaze SkipNull()] invalid char, expected 'null'")
	}
	d.pos++
	return nil
}

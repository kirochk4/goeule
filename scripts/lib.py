def cover_string(string, width, char, space):
	if len(string) + space * 2 >= width:
		return string

	left = ((width - len(string)) - space * 2) // 2
	right = left + ((width - len(string)) - space * 2)%2

	return "{}{}{}{}{}".format(
		char * left, " " * space,
		string,
		" " * space, char * right,
    )

def short_string(string, length):
	if length < len(string):
		return string[:length]
	return string

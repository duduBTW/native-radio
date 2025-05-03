package components

import (
	"github.com/dudubtw/osu-radio-native/lib"
	"github.com/dudubtw/osu-radio-native/theme"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TabsItemProps struct {
	Icon  IconName
	Label string
	Value string
}

type TabsProps struct {
	Items []TabsItemProps
	Rect  rl.Rectangle
	Value string
}

var containerPadding float32 = 5

func (c *Components) Tabs(props TabsProps) string {
	padding := lib.Padding{}
	padding.All(containerPadding)
	container := lib.NewLayout(lib.Layout{
		Direction: lib.DIRECTION_ROW,
		Padding:   padding,
		Gap:       2,
	}, props.Rect)

	var renders []func() string
	for _, tab := range props.Items {
		var drawItem func() string
		container.Render(TabItem(tab, tab.Value == props.Value, c, &drawItem))
		renders = append(renders, drawItem)
	}

	DrawRectangleRoundedPixels(
		rl.NewRectangle(props.Rect.X, props.Rect.Y, container.Size.Width, container.Size.Height),
		ROUNDED,
		rl.NewColor(13, 13, 13, 242),
	)

	newValue := props.Value
	for _, render := range renders {
		clickedValue := render()
		if clickedValue != "" {
			newValue = clickedValue
		}
	}
	return newValue
}

func TabItem(props TabsItemProps, isSelected bool, c *Components, drawItem *func() string) lib.Component {
	return func(position lib.Position, next lib.Next) {
		padding := lib.Padding{}
		// padding.Axis(0, 0)
		padding.Axis(8, 4)
		container := lib.NewLayout(lib.Layout{
			Direction: lib.DIRECTION_ROW,
			Padding:   padding,
			Gap:       8,
		}, position.ToRect(position.Contrains.Width, position.Contrains.Height))

		icon, drawIcon := TabItemIcon(props.Icon, c)
		text, drawText := TabItemText(props.Label, c.typographyMap)
		container.Render(icon)
		container.Render(text)

		rect := position.ToRect(container.Size.Width, container.Size.Height)

		getButtonStyle := func(color rl.Color) ButtonStyle {
			return ButtonStyle{Color: color, BorderRadius: ROUNDED - containerPadding}
		}

		tabStyles := ButtonStyles{
			STATE_INITIAL: getButtonStyle(rl.Fade(rl.Black, 0)),
			STATE_HOT:     getButtonStyle(rl.Fade(rl.White, 0.1)),
			STATE_ACTIVE:  getButtonStyle(rl.Fade(rl.White, 0.2)),
		}

		if isSelected {
			tabStyles = ButtonStyles{
				STATE_INITIAL: getButtonStyle(rl.DarkPurple),
				STATE_HOT:     getButtonStyle(rl.DarkPurple),
				STATE_ACTIVE:  getButtonStyle(rl.DarkPurple),
			}
		}

		*drawItem = func() string {
			isClicked := c.Button("tab-item"+props.Value, rect, tabStyles)

			drawIcon()
			drawText()
			if isClicked {
				return props.Value
			}

			return ""
		}

		next(rect)
	}
}

func TabItemIcon(name IconName, c *Components) (lib.Component, func()) {
	var p lib.Position
	return func(position lib.Position, next lib.Next) {
			p = position
			next(rl.NewRectangle(0, 0, ICON_SIZE, ICON_SIZE))
		}, func() {
			DrawIcon(name, rl.NewVector2(p.X, p.Y), &c.textures)
		}
}
func TabItemText(value string, typographyMap lib.TypographyMap) (lib.Component, func()) {
	var p lib.Position
	var themeFontSize = theme.FontSize.Regular
	var fontSize int32 = int32(themeFontSize)
	return func(position lib.Position, next lib.Next) {
			p = position
			next(position.ToRect(float32(rl.MeasureText(value, fontSize)), float32(fontSize)))
		}, func() {
			pos := rl.NewVector2(p.X, p.Y+((ICON_SIZE-float32(fontSize))/2))
			lib.Typography(value, pos, themeFontSize, theme.FontWeight.Regular, theme.Colors.Text, typographyMap)
		}
}

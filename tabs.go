package main

import (
	c "github.com/dudubtw/osu-radio-native/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TabsItemProps struct {
	Icon  c.IconName
	Label string
	Value string
}

type TabsProps struct {
	Items []TabsItemProps
	Rect  rl.Rectangle
	Value string
}

var containerPadding float32 = 5

func Tabs(props TabsProps) string {
	padding := Padding{}
	padding.All(containerPadding)
	container := NewLayout(Layout{
		Direction: DIRECTION_ROW,
		Padding:   padding,
		Gap:       2,
	}, props.Rect)

	var renders []func() string
	for _, tab := range props.Items {
		var drawItem func() string
		container.Render(TabItem(tab, tab.Value == props.Value, &drawItem))
		renders = append(renders, drawItem)
	}

	c.DrawRectangleRoundedPixels(
		rl.NewRectangle(props.Rect.X, props.Rect.Y, container.Size.Width, container.Size.Height),
		c.ROUNDED,
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

func TabItem(props TabsItemProps, isSelected bool, drawItem *func() string) Component {
	return func(position Position, next Next) {
		padding := Padding{}
		// padding.Axis(0, 0)
		padding.Axis(8, 4)
		container := NewLayout(Layout{
			Direction: DIRECTION_ROW,
			Padding:   padding,
			Gap:       8,
		}, position.ToRect(position.Contrains.Width, position.Contrains.Height))

		icon, drawIcon := TabItemIcon(props.Icon)
		text, drawText := TabItemText(props.Label)
		container.Render(icon)
		container.Render(text)

		rect := position.ToRect(container.Size.Width, container.Size.Height)

		getButtonStyle := func(color rl.Color) ButtonStyle {
			return ButtonStyle{Color: color, BorderRadius: c.ROUNDED - containerPadding}
		}

		tabStyles := ButtonStyles{
			c.STATE_INITIAL: getButtonStyle(rl.Fade(rl.Black, 0)),
			c.STATE_HOT:     getButtonStyle(rl.Fade(rl.White, 0.1)),
			c.STATE_ACTIVE:  getButtonStyle(rl.Fade(rl.White, 0.2)),
		}

		if isSelected {
			tabStyles = ButtonStyles{
				c.STATE_INITIAL: getButtonStyle(rl.DarkPurple),
				c.STATE_HOT:     getButtonStyle(rl.DarkPurple),
				c.STATE_ACTIVE:  getButtonStyle(rl.DarkPurple),
			}
		}

		*drawItem = func() string {
			isClicked := Button("tab-item"+props.Value, rect, tabStyles)

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

func TabItemIcon(name c.IconName) (Component, func()) {
	var p Position
	return func(position Position, next Next) {
			p = position
			next(rl.NewRectangle(0, 0, c.ICON_SIZE, c.ICON_SIZE))
		}, func() {
			c.DrawIcon(name, rl.NewVector2(p.X, p.Y), &textures)
		}
}
func TabItemText(value string) (Component, func()) {
	var p Position
	var fontSize int32 = 14
	return func(position Position, next Next) {
			p = position
			next(position.ToRect(float32(rl.MeasureText(value, fontSize)), float32(fontSize)))
		}, func() {
			rl.DrawText(value, int32(p.X), int32(p.Y+((c.ICON_SIZE-float32(fontSize))/2)), fontSize, rl.White)
		}
}

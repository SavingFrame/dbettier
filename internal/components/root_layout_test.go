package components

import "testing"

func TestCalculateRootLayoutSplitsExactly(t *testing.T) {
	layout := calculateRootLayout(200, 60)

	if got, want := layout.leftWidth+layout.rightWidth, 200; got != want {
		t.Fatalf("width split mismatch: got %d, want %d", got, want)
	}

	if got, want := layout.tabBarHeight+layout.tableHeight+layout.bottomHeight, layout.contentHeight; got != want {
		t.Fatalf("height split mismatch: got %d, want %d", got, want)
	}

	if got, want := layout.sqlWidth+layout.logWidth, layout.rightWidth; got != want {
		t.Fatalf("bottom-row width split mismatch: got %d, want %d", got, want)
	}

	if got, want := layout.contentHeight+layout.helpHeight+layout.statusHeight, 60; got != want {
		t.Fatalf("screen height accounting mismatch: got %d, want %d", got, want)
	}
}

func TestCalculateRootLayoutNeverNegative(t *testing.T) {
	for w := 1; w <= 200; w++ {
		for h := 1; h <= 80; h++ {
			layout := calculateRootLayout(w, h)

			if layout.leftWidth < 0 || layout.rightWidth < 0 || layout.tabBarHeight < 0 ||
				layout.tableHeight < 0 || layout.bottomHeight < 0 || layout.sqlWidth < 0 ||
				layout.logWidth < 0 || layout.contentHeight < 0 {
				t.Fatalf("negative layout values for w=%d h=%d: %+v", w, h, layout)
			}

			if layout.leftWidth+layout.rightWidth != w {
				t.Fatalf("width split mismatch for w=%d h=%d: %+v", w, h, layout)
			}

			if layout.tabBarHeight+layout.tableHeight+layout.bottomHeight != layout.contentHeight {
				t.Fatalf("content-height split mismatch for w=%d h=%d: %+v", w, h, layout)
			}

			if layout.sqlWidth+layout.logWidth != layout.rightWidth {
				t.Fatalf("sql/log split mismatch for w=%d h=%d: %+v", w, h, layout)
			}
		}
	}
}

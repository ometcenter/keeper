/* Package config определяет структуры конфигурации и некоторые интерфейсы
 */

package config

import (
	"testing"
	"time"
)

func TestServiceConfig_InitTimezone(t *testing.T) {
	time.Local, _ = time.LoadLocation("UTC")
	tests := []struct {
		name string
		s    ServiceConfig
		want string
	}{
		{
			name: "first case",
			s:    *New(),
			want: "[INFO] текущая таймзона UTC\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.InitTimezone(); got != tt.want {
				t.Errorf("ServiceConfig.InitTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceConfig_InitTimezone_MSK(t *testing.T) {
	time.Local, _ = time.LoadLocation("Europe/Moscow")
	tests := []struct {
		name string
		s    ServiceConfig
		want string
	}{
		{
			name: "first case",
			s:    *New(),
			want: "[INFO] текущая таймзона MSK\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.InitTimezone(); got != tt.want {
				t.Errorf("ServiceConfig.InitTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSettingsFromConsul(t *testing.T) {

	case1 := New()
	case1.ConsulServerAddres = "10.128.1.93:8500"

	cases := []struct {
		name string
		s    *ServiceConfig
		want error
	}{
		{
			name: "first case",
			s:    case1,
			want: nil,
		},
		{
			name: "second case",
			s:    New(),
			want: nil,
		},
	}

	for _, tt := range cases {
		// t.Run(tt.name, func(t *testing.T) {
		// 	if got := tt.s.GetSettingsFromConsul(); got != tt.want {
		// 		t.Errorf("GetSettingsFromConsul() = %v, want %v", got, tt.want)
		// 	}
		// })

		err := tt.s.GetSettingsFromConsul()
		if err != nil {
			t.Fatalf("Case %s \n GetSettingsFromConsul() returned error %q. "+
				"Error not expected.", tt.name, err)
		}
	}

}

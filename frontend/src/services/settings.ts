import {
  GetShopProfile,
  SaveShopProfile,
  GetSetting,
  SetSetting,
} from "../../wailsjs/go/app/SettingsHandler";
import { domain } from "../../wailsjs/go/models";

export type ShopProfile = domain.ShopProfile;

export const settingsApi = {
  getProfile: () => GetShopProfile(),
  saveProfile: (p: Partial<ShopProfile>) => SaveShopProfile(domain.ShopProfile.createFrom(p)),
  getSetting: (key: string) => GetSetting(key),
  setSetting: (key: string, value: string) => SetSetting(key, value),
};

export const SETTING_PRINTER_NAME = "printer_name";

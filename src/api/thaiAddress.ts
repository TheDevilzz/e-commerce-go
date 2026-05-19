export type ThaiAddressRow = {
  id: number;
  provinceCode: number;
  provinceNameEn: string;
  provinceNameTh: string;
  districtCode: number;
  districtNameEn: string;
  districtNameTh: string;
  subdistrictCode: number;
  subdistrictNameEn: string;
  subdistrictNameTh: string;
  postalCode: number;
};

export type ThaiAddressData = {
  provinces: ThaiAddressOption[];
  rows: ThaiAddressRow[];
};

type GeoThProvince = {
  id: number;
  nameTh: string;
  nameEn: string;
  districts: Array<{
    id: number;
    nameTh: string;
    nameEn: string;
  }>;
};

type GeoThDistrict = {
  id: number;
  nameTh: string;
  nameEn: string;
  provinceId: number;
  subdistricts: Array<{
    id: number;
    zipCode: number;
    nameTh: string;
    nameEn: string;
  }>;
};

export type ThaiAddressOption = {
  code: number;
  name: string;
};

const GEOTH_BASE_URL = "https://geoth.thiti.dev/api";

let thaiAddressCache: Promise<ThaiAddressData> | null = null;

export const getThaiAddressData = async () => {
  thaiAddressCache ??= Promise.all([
    fetch(`${GEOTH_BASE_URL}/provinces-with-districts/all`),
    fetch(`${GEOTH_BASE_URL}/districts-with-subdistricts/all`),
  ]).then(async ([provincesResponse, districtsResponse]) => {
    if (!provincesResponse.ok || !districtsResponse.ok) {
      throw new Error("Unable to load GeoTH address data");
    }

    const provinces = (await provincesResponse.json()) as GeoThProvince[];
    const districts = (await districtsResponse.json()) as GeoThDistrict[];
    const provincesById = new Map(provinces.map((province) => [province.id, province]));
    const rows = districts.flatMap((district) => {
      const province = provincesById.get(district.provinceId);

      if (!province) {
        return [];
      }

      return district.subdistricts.map((subdistrict) => ({
        id: subdistrict.id,
        provinceCode: province.id,
        provinceNameEn: province.nameEn,
        provinceNameTh: province.nameTh,
        districtCode: district.id,
        districtNameEn: district.nameEn,
        districtNameTh: district.nameTh,
        subdistrictCode: subdistrict.id,
        subdistrictNameEn: subdistrict.nameEn,
        subdistrictNameTh: subdistrict.nameTh,
        postalCode: subdistrict.zipCode,
      }));
    });

    return {
      provinces: provinces
        .map((province) => ({ code: province.id, name: province.nameTh }))
        .sort((a, b) => a.name.localeCompare(b.name, "th")),
      rows,
    };
  });

  return thaiAddressCache;
};

export const getThaiProvinces = (data: ThaiAddressData | null): ThaiAddressOption[] => data?.provinces ?? [];

export const getThaiDistricts = (data: ThaiAddressData | null, provinceCode: number): ThaiAddressOption[] => {
  const districts = new Map<number, string>();

  (data?.rows ?? [])
    .filter((row) => row.provinceCode === provinceCode)
    .forEach((row) => {
      districts.set(row.districtCode, row.districtNameTh);
    });

  return Array.from(districts, ([code, name]) => ({ code, name })).sort((a, b) =>
    a.name.localeCompare(b.name, "th"),
  );
};

export const getThaiSubdistricts = (data: ThaiAddressData | null, districtCode: number): ThaiAddressOption[] =>
  (data?.rows ?? [])
    .filter((row) => row.districtCode === districtCode)
    .map((row) => ({ code: row.subdistrictCode, name: row.subdistrictNameTh }))
    .sort((a, b) => a.name.localeCompare(b.name, "th"));

export const findThaiAddressRow = (data: ThaiAddressData | null, subdistrictCode: number) =>
  data?.rows.find((row) => row.subdistrictCode === subdistrictCode);

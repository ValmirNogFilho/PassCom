import giro_brand from "../assets/brand_giro.svg"
import rumos_brand from "../assets/brand_rumos.svg"
import boreal_brand from "../assets/brand_boreal.svg"

export const findCompany = (company) => {
    switch (company) {
      case "giro":
        return giro_brand
      case "boreal":
        return boreal_brand
      default:
        return rumos_brand
    }
  }
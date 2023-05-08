import { Component } from '@angular/core';
import {HttpClient} from "@angular/common/http";

const backendUrl = "http://localhost:8000/api/order"

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'ng-app';


  constructor(private httpClient: HttpClient) {
  }

  order_uuid: string = '';

  data_json: string = '';

  showResults = false;

  search() {
    let url = this.order_uuid == '' ? backendUrl : `${backendUrl}/${this.order_uuid}`
    this.httpClient.get(url).subscribe(
      (data: any) => {
        this.data_json = data
        this.showResults = true
      },
      (error) => {
        this.data_json = error
        this.showResults = true
      }
    )
  }


}

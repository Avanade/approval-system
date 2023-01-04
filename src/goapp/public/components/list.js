const list = ({
    callback,
    renderItem,
  }) => {
    return {
    // FILTER
    search : '',
    filter : 10,
    page : 0,

    // DISPLAY
    total : 0,
    showStart : 0,
    showEnd : 0,
    isLoading : false,

    items : [],
    async init(){
        await this.load();
    },
    async load(){
        this.isLoading = true
        this.showStart = 0
        this.showEnd = 0
        const {data, total} = await callback({
            filter : this.filter,
            page : this.page,
            search : this.search,
        })

        this.items = data
        this.total = total
        
        this.isLoading = false

        if (this.items == null || this.items.length == 0) return;

        this.showStart = this.items.length > 0 ? ((this.page * this.filter) + 1) : 0;
        this.showEnd = (this.page * this.filter) + this.items.length;
    },
    //EVENT HANDLERS
    onChangeFilterHandler(e){
        this.page = 0,
        this.total = 0,
        this.filter = parseInt(e.target.value);
        this.load()
    },
    onSearchSubmitHandler(e){
        this.page = 0,
        this.total = 0,
        this.search = e.target.value;
        this.load();
    },
    onNextPageHandler(){
        if (!this.nextPageEnabled()) return;

        this.page = this.page + 1

        this.load()
    },
    onPreviousPageHandler(){
        if (!this.previousPageEnabled()) return;

        this.page = this.page - 1

        this.load();
    },
    //FUNCTIONS
    nextPageEnabled(){
        return this.page < Math.ceil(this.total/this.filter) - 1
    },
    previousPageEnabled(){
        return this.page > 0
    },
    //RENDER
    render(item){
        return renderItem(item)
    },
    template : `<nav class="bg-white flex items-center justify-between" aria-label="header">
            <div class="sm:block">
                <div>
                    <label for="filter" class="block text-sm font-medium text-gray-700">Filter</label>
                    <select @change="onChangeFilterHandler" x-model="filter" id="filter" name="filter" class="mt-1 block w-20 pl-3 pr-10 py-2 text-base text-center border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md">
                        <option>5</option>
                        <option>10</option>
                        <option>20</option>
                        <option>50</option>
                        <option>100</option>
                    </select>
                </div>
            </div>
            <div class="flex justify-between sm:justify-end">
                <div class="sm:col-span-3">
                    <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
                    <div class="mt.-1">
                        <input @keyup.enter="onSearchSubmitHandler" type="text" name="search" id="search" class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md" x-model="search">
                    </div>
                </div>
            </div>
        </nav>

        <div x-show='isLoading' x-transition>
            <svg 
                role="status" 
                class="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-[#FF5800] m-auto my-5"
                viewBox="0 0 100 101" 
                fill="none" 
                xmlns="http://www.w3.org/2000/svg">
                <path
                d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                fill="currentColor" />
                <path
                d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                fill="currentFill" />
            </svg>
        </div>
        
        <div x-show="!isLoading" x-transition>
            <template x-if="!items">
                <p class="text-center my-5">NO RESULT FOUND</p>
            </template>
            <template x-if="items.length != 0">
                <ul role="list" class="divide-y divide-gray-300 my-3">
                    <template x-for="item in items">
                        <li x-html="render(item)">
                        </li>
                    </template>
                </ul>
            </template>
        </div>

        <nav class="bg-white py-3 flex items-center justify-between border-t border-gray-200" aria-label="Pagination">
            <div class="sm:block">
                <p class="text-sm text-gray-700">
                    Showing
                    <span class="font-medium" x-text="showStart"></span>
                    to
                    <span class="font-medium" x-text="showEnd"></span>
                    of
                    <span class="font-medium" x-text="total"></span>
                    results
                </p>
            </div>
            <div class="flex justify-between sm:justify-end">
                <button x-bind:disabled="!previousPageEnabled()" x-on:click="onPreviousPageHandler" href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-200"> Previous </button>
                <button x-bind:disabled="!nextPageEnabled()" x-on:click="onNextPageHandler" href="#" class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-200"> Next </button>
            </div>
        </nav>`
    }
  }
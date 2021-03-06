<template>
    <v-dialog v-model="show" persistent max-width="350" class="p-photo-album-dialog" @keydown.esc="cancel">
        <v-card raised elevation="24">
            <v-container fluid class="pb-2 pr-2 pl-2">
                <v-layout row wrap>
                    <v-flex xs3 text-xs-center>
                        <v-icon size="54" color="grey lighten-1" v-if="newAlbum">create_new_folder</v-icon>
                        <v-icon size="54" color="grey lighten-1" v-else>folder</v-icon>
                    </v-flex>
                    <v-flex xs9 text-xs-left align-self-center>
                        <v-autocomplete
                                v-model="album"
                                browser-autocomplete="off"
                                hint="Album Name"
                                :items="items"
                                :search-input.sync="search"
                                :loading="loading"
                                hide-details
                                item-text="AlbumName"
                                item-value="AlbumUUID"
                                :label="labels.select"
                                color="secondary-dark"
                                flat solo
                        >
                        </v-autocomplete>
                    </v-flex>
                    <v-flex xs12 text-xs-right class="pt-3">
                        <v-btn @click.stop="cancel" depressed color="grey lighten-3" class="p-photo-dialog-cancel">
                            <translate>Cancel</translate>
                        </v-btn>
                        <v-btn color="blue-grey lighten-2" depressed dark @click.stop="confirm"
                               class="p-photo-dialog-confirm">
                            <span v-if="newAlbum">{{ labels.createAlbum }}</span>
                            <span v-else>{{ labels.addToAlbum }}</span>
                        </v-btn>
                    </v-flex>
                </v-layout>
            </v-container>
        </v-card>
    </v-dialog>
</template>
<script>
    import Album from "../model/album";

    export default {
        name: 'p-photo-album-dialog',
        props: {
            show: Boolean,
        },
        data() {
            return {
                loading: false,
                search: null,
                newAlbum: null,
                album: "",
                albums: [],
                items: [],
                labels: {
                    select: this.$gettext("Select album"),
                    addToAlbum: this.$gettext("Add to album"),
                    createAlbum: this.$gettext("Create album"),
                }
            }
        },
        methods: {
            cancel() {
                this.$emit('cancel');
            },
            confirm() {
                if(this.album === "" && this.newAlbum) {
                    console.log("NEW", this.album, this.newAlbum);
                    this.loading = true;

                    this.newAlbum.save().then((a) => {
                        this.loading = false;
                        this.$emit('confirm', a.AlbumUUID);
                    });
                } else {
                    console.log("OLD", this.album, this.newAlbum);
                    this.$emit('confirm', this.album);
                }
            },
            queryServer(q) {
                if(this.loading) {
                    return;
                }

                this.loading = true;

                const params = {
                    q: q,
                    count: 1000,
                    offset: 0,
                };

                Album.search(params).then(response => {
                    this.loading = false;

                    if(response.models.length > 0 && !this.album) {
                        this.album = response.models[0].AlbumUUID;
                    }

                    this.albums = response.models;
                    this.items = [...this.albums];
                }).catch(() => this.loading = false);
            },
        },
        watch: {
            search (q) {
                const exists = this.albums.findIndex((album) => album.AlbumName === q);

                if (exists !== -1 || !q) {
                    this.items = this.albums;
                    this.newAlbum = null;
                } else {
                    this.newAlbum = new Album({AlbumName: q, AlbumUUID: ""});
                    this.items = this.albums.concat([this.newAlbum]);
                }
            },
        },
        created() {
            this.queryServer("");
        },
    }
</script>
